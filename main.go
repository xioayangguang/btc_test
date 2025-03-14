package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"log"
	"os"
)

var client = &rpcclient.Client{}

func init() {
	// 使用tls链接，所以需要导入btcd生成的rpc证书
	cert, err := os.ReadFile("./rpc.cert")
	if err != nil {
		panic(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host: "192.168.31.4:8334",
		//Host:         "127.0.0.1:8334",
		User:         "root",
		Pass:         "root",
		HTTPPostMode: true,
		//DisableTLS:   true,
		Certificates: cert,
	}
	//var err error
	client, err = rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	hash, err := client.GetBlockHash(blockCount)
	if err != nil {
		log.Fatal(err)
	}
	block, err := client.GetBlock(hash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Block version: %v\n", block.Header.Version)
	fmt.Printf("Block hash: %v\n", block.BlockHash())
	fmt.Printf("Block previous hash: %v\n", block.Header.PrevBlock)
	fmt.Printf("Block merkle root: %v\n", block.Header.MerkleRoot)
	fmt.Printf("Block timestamp: %v\n", block.Header.Timestamp)
	fmt.Printf("Block bits: %v\n", block.Header.Bits)
	fmt.Printf("Block nonce: %v\n", block.Header.Nonce)
	fmt.Printf("Number of transactions in block: %v\n", len(block.Transactions))
}

func GenerateBTCTest() (string, string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", err
	}
	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, false)
	if err != nil {
		return "", "", err
	}
	pubKeySerial := privKey.PubKey().SerializeUncompressed()
	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.TestNet3Params)
	if err != nil {
		return "", "", err
	}
	return privKeyWif.String(), pubKeyAddress.EncodeAddress(), nil
}

// BuildTxOut 构建一个比特币交易输出（TxOut）
func BuildTxOut(addr string, amount int64, params chaincfg.Params) (*wire.TxOut, []byte, error) {
	// 解析比特币地址
	destinationAddress, err := btcutil.DecodeAddress(addr, &params)
	if err != nil {
		return nil, nil, err
	}

	// 生成支付到地址的脚本
	pkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return nil, nil, err
	}

	// 创建一个新的交易输出，金额单位为 satoshis
	return wire.NewTxOut(amount, pkScript), pkScript, nil
}

// GetUTXOs 获取指定比特币地址的所有未花费交易输出（UTXOs）
func GetUTXOs(addr string) ([]*btcjson.ListUnspentResult, error) {
	// 解析比特币地址
	address, err := btcutil.DecodeAddress(addr, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	// 使用SearchRawTransactionsVerbose获取与地址相关的所有交易
	transactions, err := client.SearchRawTransactionsVerbose(address, 0, 100, true, false, nil)
	if err != nil {
		return nil, err
	}

	// 用于存储UTXO的切片
	utxos := []*btcjson.ListUnspentResult{}

	// 遍历所有交易
	for _, tx := range transactions {
		// 将交易ID字符串转换为链哈希对象
		txid, err := chainhash.NewHashFromStr(tx.Txid)
		if err != nil {
			log.Fatalf("Invalid txid: %v", err)
		}

		// 遍历交易的输出
		for _, vout := range tx.Vout {
			// 检查输出地址是否是我们关心的地址
			if vout.ScriptPubKey.Address != addr {
				continue
			}

			// 使用GetTxOut方法获取交易输出，确认该输出是否未花费
			utxo, err := client.GetTxOut(txid, vout.N, true)
			if err != nil {
				panic(err)
			}

			// 如果交易输出未花费，则将其添加到UTXO切片中
			if utxo != nil {
				utxo := &btcjson.ListUnspentResult{
					TxID:          tx.Txid,
					Vout:          uint32(vout.N),
					Address:       addr,
					ScriptPubKey:  vout.ScriptPubKey.Hex,
					Amount:        vout.Value, // 单位为BTC
					Confirmations: int64(tx.Confirmations),
					Spendable:     true,
				}
				utxos = append(utxos, utxo)
			}
		}
	}

	// 返回UTXO集合
	return utxos, nil
}

func BuildTxIn(wif *btcutil.WIF, amount int64, txOut *wire.TxOut, params *chaincfg.Params) (*wire.MsgTx, error) {
	// 解析比特币地址
	//fromAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), params)
	fromAddr, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), params)
	if err != nil {
		return nil, errors.Wrap(err, "解析比特币地址失败")
	}

	// 获取UTXOs
	utxos, err := GetUTXOs(fromAddr.EncodeAddress())
	if err != nil {
		return nil, errors.Wrap(err, "获取UTXOs失败")
	}

	msgTx := wire.NewMsgTx(wire.TxVersion)
	// 创建一个新的交易输入，金额单位为 satoshis
	totalInput := int64(0)
	for _, utxo := range utxos {
		// totalInput 大于 amount，用于计算交易费
		if totalInput > amount {
			break
		}
		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, errors.Wrap(err, "解析交易哈希失败")
		}
		//pkScript, _ := txscript.PayToAddrScript(addr)
		//outputs = append(outputs, wire.NewTxOut(leftToMe, pkScript))
		//prevOutputScript, err := hex.DecodeString(utxo.ScriptPubKey)
		//if err != nil {
		//	panic(err)
		//}
		outPoint := &wire.OutPoint{Hash: *txHash, Index: utxo.Vout}
		txIn := wire.NewTxIn(outPoint, nil, nil)
		msgTx.AddTxIn(txIn)
		totalInput += int64(utxo.Amount * 1e8)
	}
	msgTx.AddTxOut(txOut)

	// 交易费
	// 假定交易费率为每字节 10sat
	fee := int64(msgTx.SerializeSize()) * 10
	// 找零
	change := totalInput - amount
	// 这里假定找零一定大于交易费，交易费太少的话可能导致交易一直无法确认
	// 如果change <= fee的话，零钱会转给出块的矿工
	if change > fee {
		changePkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			return nil, errors.Wrap(err, "生成找零地址的脚本失败")
		}
		txOut := wire.NewTxOut(change-fee, changePkScript)
		msgTx.AddTxOut(txOut)
	}

	// 签署交易
	// 发送方地址为SegWit的P2WPKH 地址，所以要消费该地址的UTXO，只能通过见证输入进行消费
	for i, _ := range msgTx.TxIn {
		prevOutputScript, err := hex.DecodeString(utxos[i].ScriptPubKey)
		if err != nil {
			panic(err)
		}
		//pkScript := prevPkScripts[i]
		//var script []byte
		script, err := txscript.SignatureScript(msgTx, int(utxos[i].Vout), prevOutputScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			return nil, errors.Wrap(err, "解析交易哈希失败")
		}
		msgTx.TxIn[i].SignatureScript = script

		//txHash, err := chainhash.NewHashFromStr(utxos[i].TxID)
		//if err != nil {
		//	return nil, errors.Wrap(err, "解析交易哈希失败")
		//}
		//outPoint := wire.OutPoint{Hash: *txHash, Index: uint32(utxos[i].Vout)}
		//prevOutputFetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		//	outPoint: {Value: int64(utxos[i].Amount * 1e8), PkScript: prevOutputScript},
		//})
		//sigHashes := txscript.NewTxSigHashes(msgTx, prevOutputFetcher)
		//sigScript, err := txscript.WitnessSignature(msgTx, sigHashes, int(utxos[i].Vout), int64(utxos[i].Amount*1e8), prevOutputScript, txscript.SigHashAll, wif.PrivKey, true)
		//sigScript, err := txscript.SignatureScript(msgTx, int(utxos[i].Vout), prevOutputScript, txscript.SigHashAll, wif.PrivKey, true)
		//if err != nil {
		//	return nil, errors.Wrap(err, "签名交易失败")
		//}
		//txIn.Witness = sigScript
	}
	return msgTx, nil
}

//
//func generateLockingScript(pubKey []byte) []byte {
//	// 生成公钥哈希
//	pubKeyHash := btcutil.Hash160(pubKey)
//	// 生成P2PKH锁定脚本
//	scriptPubKey, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).
//		AddOp(txscript.OP_HASH160).AddData(pubKeyHash).
//		AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).Script()
//	if err != nil {
//		panic(err)
//	}
//	return scriptPubKey
//}
//
//func createUnlockingScript(privKey *btcutil.WIF, pubKey, prevScriptPubKey []byte, tx *wire.MsgTx) []byte {
//	// 创建签名
//	sig, err := txscript.RawTxInSignature(tx, 0, prevScriptPubKey, txscript.SigHashAll, privKey.PrivKey)
//	if err != nil {
//		panic(err)
//	}
//	// 生成解锁脚本
//	scriptSig, err := txscript.NewScriptBuilder().AddData(sig).AddData(pubKey).Script()
//	if err != nil {
//		panic(err)
//	}
//	return scriptSig
//}
//
//func main() {
//	privKey, pubKey := generateKeys()
//	fmt.Println("私钥:", privKey)
//	fmt.Println("公钥:", pubKey)
//	scriptPubKey := generateLockingScript(pubKey)
//	fmt.Println("锁定脚本:", scriptPubKey)
//	// 示例前序交易ID和输出索引
//	prevTxHash, _ := chainhash.NewHashFromStr("previous_txid")
//	prevTxOut := wire.NewTxOut(0, scriptPubKey)
//	// 创建新交易
//	tx := wire.NewMsgTx(wire.TxVersion)
//	tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(prevTxHash, 0), nil, nil))
//	tx.AddTxOut(prevTxOut)
//	scriptSig := createUnlockingScript(privKey, pubKey, scriptPubKey, tx)
//	tx.TxIn[0].SignatureScript = scriptSig
//	// 序列化交易
//	buf := &bytes.Buffer{}
//	err := tx.Serialize(buf)
//	if err != nil {
//		panic(err)
//	}
//	// 将交易转换为十六进制字符串
//	txHex := fmt.Sprintf("%x", hex.EncodeToString(buf.Bytes()))
//	fmt.Println("交易:", txHex)
//	// 在此处添加广播交易的代码，通常需要连接到比特币节点或使用第三方API
//}
