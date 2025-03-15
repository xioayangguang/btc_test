package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"log"
	"testing"
)

//
//var client *rpcclient.Client
//
//func init() {
//	cert, err := os.ReadFile("./btcd/rpc.cert")
//	if err != nil {
//		panic(err)
//	}
//	//Connect to local bitcoin core RPC server using HTTP POST mode.
//	connCfg := &rpcclient.ConnConfig{
//		Host:         "localhost:18332",
//		User:         "your_rpc_user",
//		Pass:         "your_rpc_password",
//		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
//		DisableTLS:   true, // Bitcoin core does not provide TLS by default
//		Certificates: cert,
//	}
//	c, err := rpcclient.New(connCfg, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	client = c
//}

// 研究https://mempool.space/zh/tx/fc1e70787d2bc345eb9cf0e3bc55eeeaa6eb3bb3365e0995d6a62e4dfabf8caf这个地址的交易类型
// 通过私钥生成所有的私人地址，多签名钱包不属于私人地址。因为和其他不同的钱包地址组合都是不同的多签钱包地址
func NewBTCAddress() {
	var pubKey *secp256k1.PublicKey
	if true {
		//https://key.tokenpocket.pro/?locale=zh#/?network=BTC
		wifStr := "L5A7ZqEJswd2RXFdDeWYj2kmLZKnC2HWvzDosdZ8TxsFXvnSrcyb"
		// 解析 WIF 格式的私钥
		wif, err := btcutil.DecodeWIF(wifStr)
		if err != nil {
			log.Fatalf("解析 WIF 失败: %v", err)
		}
		pubKey = wif.PrivKey.PubKey()
	} else {
		privKey, err := btcec.NewPrivateKey()
		if err != nil {
			panic(err)
		}
		fmt.Println("===================Private Key===============================")
		fmt.Printf("Binary Private Key: %x\n", privKey.Serialize())
		WIFPrivate, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("compressed WIF Private Key: %s\n", WIFPrivate.String())
		wifPrivate, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Uncompressed WIF Private Key: %s\n", wifPrivate.String())
		fmt.Println()
		fmt.Println("=================== Public Key===============================")
		// 导出公钥
		pubKey = privKey.PubKey()
	}

	compressedPubKey := pubKey.SerializeCompressed()
	fmt.Printf("compressed  Public Key: %x\n", compressedPubKey)
	uncompressedPubKey := pubKey.SerializeUncompressed()
	fmt.Printf("Uncompressed Public Key: %x\n", uncompressedPubKey)
	pubKeyHashUnCompressed := btcutil.Hash160(uncompressedPubKey)
	pubKeyHashCompressed := btcutil.Hash160(compressedPubKey)
	fmt.Println()
	fmt.Println("===================P2PKH===============================")
	addressPKH, err := btcutil.NewAddressPubKey(compressedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("compressed P2PKH Address: %s\n", addressPKH.EncodeAddress())
	// 生成BTC地址（P2PKH）
	addressPKH, err = btcutil.NewAddressPubKey(uncompressedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("unCompressed P2PKH Address: %s\n", addressPKH.EncodeAddress())
	fmt.Println()
	fmt.Println("===================Bech32编码地址是专为SegWit 原生隔离见证地址===================")
	// 生成Bech32地址（P2WPKH）
	addressBech32_1, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHashCompressed, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("compressed Bech32 Address: %s\n", addressBech32_1.EncodeAddress())
	addressBech32_2, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHashUnCompressed, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("unCompressed Bech32 Address: %s\n", addressBech32_2.EncodeAddress())
	fmt.Println()
	fmt.Println("===================P2SH-P2WPKH  兼容隔离见证地址===============================")
	fmt.Println("3开头的地址一部分是多签地址(背后被其他多个其他类型的地址操控)，一部分是隔离见证兼容地址（私人控制，现在生成的就是私人控制类型的地址），具体属于那种地址可以去区块链浏览器查看地址类型")
	fmt.Println("3开头的多签地址 本钱包和其他任何地址组合成的多签地址都不一样")
	witnessScript, err := txscript.PayToAddrScript(addressBech32_1)
	if err != nil {
		log.Fatalf("构造 P2WPKH 脚本失败: %v", err)
	}
	address1, err := btcutil.NewAddressScriptHash(witnessScript, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("生成 P2SH 地址失败: %v", err)
	}
	fmt.Printf("compressed P2SH Address: %s\n", address1.EncodeAddress())

	witnessScript, err = txscript.PayToAddrScript(addressBech32_2)
	if err != nil {
		log.Fatalf("构造 P2WPKH 脚本失败: %v", err)
	}
	address1, err = btcutil.NewAddressScriptHash(witnessScript, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("生成 P2SH 地址失败: %v", err)
	}
	fmt.Printf("unCompressed P2SH Address: %s\n", address1.EncodeAddress())
	fmt.Println()
	fmt.Println("===================taprootAddr===============================")
	taprootAddr, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pubKey)), &chaincfg.MainNetParams)
	fmt.Printf("taprootAddr Address: %s\n", taprootAddr.EncodeAddress())
	fmt.Println()
}

func TestA(t *testing.T) {
	NewBTCAddress()

	GenP2SHAddress()
}

func GenP2SHAddress() {
	//   假设我们要创建一个2-of-3的多重签名脚本（这意味着要花费比特币，需要3个可能的签名者中的2个签名）。
	_, pubKey1 := generateKeys()
	address1Pub, err := btcutil.NewAddressPubKey(pubKey1, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error   NewAddressPubKey:", err)
		return
	}
	_, pubKey2 := generateKeys()
	address2Pub, err := btcutil.NewAddressPubKey(pubKey2, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error   NewAddressPubKey:", err)
		return
	}
	_, pubKey3 := generateKeys()
	address3Pub, err := btcutil.NewAddressPubKey(pubKey3, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error   NewAddressPubKey:", err)
		return
	}
	//   创建多重签名赎回脚本
	redeemScript, err := txscript.MultiSigScript([]*btcutil.AddressPubKey{address1Pub, address2Pub, address3Pub}, 2)
	if err != nil {
		fmt.Println("Error   creating   redeem   script:", err)
		return
	}
	// 打印赎回脚本（以16进制表示）
	fmt.Printf("赎回脚本 Redeem Script:   %x\n", redeemScript)
	//计算P2SH地址
	//redeemScriptHash := btcutil.Hash160(redeemScript)
	address, err := btcutil.NewAddressScriptHash(redeemScript, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Error   creating   P2SH   address:", err)
		return
	}
	fmt.Println("===================生成多签地址 P2SH===============================")
	//   打印P2SH地址
	fmt.Println("多签地址 P2SH Address:", address.EncodeAddress())
}

// 生成私钥和公钥
func generateKeys() (*btcutil.WIF, []byte) {
	//   生成私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
	priKey, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
	if err != nil {
		panic(err)
	}
	//   生成公钥
	pubKey := privateKey.PubKey().SerializeCompressed()
	return priKey, pubKey
}

// https://studygolang.com/articles/12303
// https://mempool.space/zh/testnet/address/mtvJM2gFAASs6yqifaym1i3pY8GzspNxQ8
// https://mempool.space/testnet/api/address/mtvJM2gFAASs6yqifaym1i3pY8GzspNxQ8/utxo
// https://api.blockcypher.com/v1/btc/test3/addrs/mtvJM2gFAASs6yqifaym1i3pY8GzspNxQ8
func TestB(t *testing.T) {
	address := "mtvJM2gFAASs6yqifaym1i3pY8GzspNxQ8"
	var balance int64 = 19916    // 余额    //todo  //替换成自己的
	var fee int64 = 0.0001 * 1e8 // 交易费
	var leftToMe = balance - fee // 余额-交易费就是剩下再给我的

	// 1. 构造输出
	var outputs []*wire.TxOut
	// 1.1 输出1, 给自己转剩下的钱
	addr, _ := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	pkScript, _ := txscript.PayToAddrScript(addr)
	outputs = append(outputs, wire.NewTxOut(leftToMe, pkScript))
	// 1.2 输出2, 添加文字
	comment := "我是比特币测试zxx, 哈哈"
	pkScript, _ = txscript.NullDataScript([]byte(comment))
	outputs = append(outputs, wire.NewTxOut(int64(0), pkScript))

	// 2. 构造输入
	//prevTxHash := "48eea09764713f3dadcfed29490ab5e288299e01e571e1f7a1396a75ce38e067" //替换成自己的
	prevTxHash := "bc78192912e9a0a4ec115626df86900dd50833d2511c1379dd09607dec66e48e" //替换成自己的
	prevTxOutputN := uint32(1)                                                       //替换成自己的
	hash, _ := chainhash.NewHashFromStr(prevTxHash)                                  // tx hash
	outPoint := wire.NewOutPoint(hash, prevTxOutputN)                                // 第几个输出
	txIn := wire.NewTxIn(outPoint, nil, nil)
	inputs := []*wire.TxIn{txIn}

	tx := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
	}

	// 3. 签名
	//prevPkScriptHex := "76a91489a7f0117eaf47d8b4af740c66116e35ffe1bea988ac" //替换成自己的
	prevPkScriptHex := "76a9149303f0b74ceadd782fee549903bd336894d07e1e88ac" //替换成自己的
	prevPkScript, _ := hex.DecodeString(prevPkScriptHex)
	prevPkScripts := make([][]byte, 1)
	prevPkScripts[0] = prevPkScript

	privKey := "92VtWgn9krfwZup8Vu2g41NXBX4gZvAwWrRfknx69MJrVkT8jhm" // 私钥   //替换成自己的
	sign(tx, privKey, prevPkScripts)

	// 4. 输出Hex
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
	}
	txHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("hex", txHex)
}

// 010000000167e038ce756a39a1f7e171e5019e2988e2b50a4929edcfad3d3f716497a0ee48000000008a473044022100c9fa5201b4ed6d889c37e9173ca108c4302d4d042d1d7a4815fa5085ead1b78f021f6795ba4d95fc90659d55f1387c028a8cb5de32cea3236a6a4069fc33aca70501410491c38af613c6731597f7bdfaa2dd6f98cc1acec71bf860509e8559f8bc5ab4dd8e8eaa6d236ee81e2156385c5d4166403c21a766f55ed99cee36a53df9dd518bffffffff02a04bde03000000000000000000000000001c6a1ae8bf99e698afe4b880e4b8aae79599e8a8802c20e59388e5938800000000
// 010000000167e038ce756a39a1f7e171e5019e2988e2b50a4929edcfad3d3f716497a0ee48000000008a473044022100c9fa5201b4ed6d889c37e9173ca108c4302d4d042d1d7a4815fa5085ead1b78f021f6795ba4d95fc90659d55f1387c028a8cb5de32cea3236a6a4069fc33aca70501410491c38af613c6731597f7bdfaa2dd6f98cc1acec71bf860509e8559f8bc5ab4dd8e8eaa6d236ee81e2156385c5d4166403c21a766f55ed99cee36a53df9dd518bffffffff02a04bde03000000000000000000000000001c6a1ae8bf99e698afe4b880e4b8aae79599e8a8802c20e59388e5938800000000
// 01000000018ee466ec7d6009dd79131c51d23308d50d9086df265611eca4a0e912291978bc010000008b483045022100ca84da4103edaf5de3c357679216dc3ce81fb2900c840a5de5280afe6c66031a02203d161a65052af75f223854b2e36b1bdd434bdf6b8c3073a4b30cbfed59019dd001410491c38af613c6731597f7bdfaa2dd6f98cc1acec71bf860509e8559f8bc5ab4dd8e8eaa6d236ee81e2156385c5d4166403c21a766f55ed99cee36a53df9dd518bffffffff02bc26000000000000000000000000000000226a20e68891e698afe6af94e789b9e5b881e6b58be8af957a78782c20e59388e5938800000000
// 签名
func sign(tx *wire.MsgTx, privKeyStr string, prevPkScripts [][]byte) {
	inputs := tx.TxIn
	//众所周知，比特币常见到的私钥格式有三种，分别是16进制格式的，WIF格式，以及WIF压缩格式
	wif, err := btcutil.DecodeWIF(privKeyStr)
	fmt.Println("wif err", err)
	privKey := wif.PrivKey
	for i := range inputs {
		pkScript := prevPkScripts[i]
		var script []byte
		script, err = txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, privKey, false)
		inputs[i].SignatureScript = script
	}
}

func TestC(t *testing.T) {
	bc, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("block  count:  %d\n", bc)

	address := "mtvJM2gFAASs6yqifaym1i3pY8GzspNxQ8"

	//GetUTXOs(address)

	addr, _ := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)

	a, err := client.ListUnspentMinMaxAddresses(6, 999999999, []btcutil.Address{addr})
	fmt.Println(err)
	fmt.Println(a)

	aa, err := client.ListUnspent()
	fmt.Println(aa)

	v, err := client.BackendVersion()
	log.Printf("Block count: %d", v)

	b, err := client.GetBalance("*")
	log.Printf("Block count: %d", b)

	prevBlockHash, err := client.GetBestBlockHash()
	if err != nil {
		t.Fatalf("unable to get prior block hash: %v", err)
	}
	prevBlock, err := client.GetBlock(prevBlockHash)
	if err != nil {
		t.Fatalf("unable to get block: %v", err)
	}

	fmt.Println(prevBlock)
	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	count, err := client.GetBlockCount()

	fmt.Println(err)
	fmt.Println(count)

}

//
//// BuildTxOut 构建一个比特币交易输出（TxOut）
//func BuildTxOut(addr string, amount int64, params chaincfg.Params) (*wire.TxOut, []byte, error) {
//	// 解析比特币地址
//	destinationAddress, err := btcutil.DecodeAddress(addr, &params)
//	if err != nil {
//		return nil, nil, err
//	}
//	// 生成支付到地址的脚本
//	pkScript, err := txscript.PayToAddrScript(destinationAddress)
//	if err != nil {
//		return nil, nil, err
//	}
//	// 创建一个新的交易输出，金额单位为 satoshis
//	return wire.NewTxOut(amount, pkScript), pkScript, nil
//}
//
//// GetUTXOs 获取指定比特币地址的所有未花费交易输出（UTXOs）
//func GetUTXOs(addr string) ([]*btcjson.ListUnspentResult, error) {
//	// 解析比特币地址
//	address, err := btcutil.DecodeAddress(addr, &chaincfg.TestNet3Params)
//	if err != nil {
//		return nil, err
//	}
//
//	// 使用SearchRawTransactionsVerbose获取与地址相关的所有交易
//	transactions, err := client.SearchRawTransactionsVerbose(address, 0, 100, true, false, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	// 用于存储UTXO的切片
//	utxos := []*btcjson.ListUnspentResult{}
//
//	// 遍历所有交易
//	for _, tx := range transactions {
//		// 将交易ID字符串转换为链哈希对象
//		txid, err := chainhash.NewHashFromStr(tx.Txid)
//		if err != nil {
//			log.Fatalf("Invalid txid: %v", err)
//		}
//
//		// 遍历交易的输出
//		for _, vout := range tx.Vout {
//			// 检查输出地址是否是我们关心的地址
//			if vout.ScriptPubKey.Address != addr {
//				continue
//			}
//
//			// 使用GetTxOut方法获取交易输出，确认该输出是否未花费
//			utxo, err := client.GetTxOut(txid, vout.N, true)
//			if err != nil {
//				panic(err)
//			}
//
//			// 如果交易输出未花费，则将其添加到UTXO切片中
//			if utxo != nil {
//				utxo := &btcjson.ListUnspentResult{
//					TxID:          tx.Txid,
//					Vout:          uint32(vout.N),
//					Address:       addr,
//					ScriptPubKey:  vout.ScriptPubKey.Hex,
//					Amount:        vout.Value, // 单位为BTC
//					Confirmations: int64(tx.Confirmations),
//					Spendable:     true,
//				}
//				utxos = append(utxos, utxo)
//			}
//		}
//	}
//	// 返回UTXO集合
//	return utxos, nil
//}
//
//func BuildTxIn(wif *btcutil.WIF, amount int64, txOut *wire.TxOut, params *chaincfg.Params) (*wire.MsgTx, error) {
//	msgTx := wire.NewMsgTx(wire.TxVersion)
//	msgTx.AddTxOut(txOut)
//
//	// 解析比特币地址
//	fromAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), params)
//	if err != nil {
//		return nil, errors.Wrap(err, "解析比特币地址失败")
//	}
//	// 获取UTXOs
//	utxos, err := GetUTXOs(fromAddr.EncodeAddress())
//	if err != nil {
//		return nil, errors.Wrap(err, "获取UTXOs失败")
//	}
//
//	// 创建一个新的交易输入，金额单位为 satoshis
//	totalInput := int64(0)
//	for _, utxo := range utxos {
//		// totalInput 大于 amount，用于计算交易费
//		if totalInput > amount {
//			break
//		}
//		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
//		if err != nil {
//			return nil, errors.Wrap(err, "解析交易哈希失败")
//		}
//		txIn := wire.NewTxIn(&wire.OutPoint{
//			Hash:  *txHash,
//			Index: utxo.Vout,
//		}, nil, nil)
//
//		msgTx.AddTxIn(txIn)
//		totalInput += int64(utxo.Amount * 1e8)
//	}
//
//	// 交易费
//	// 假定交易费率为每字节 1sat
//	fee := int64(msgTx.SerializeSize())
//	// 找零
//	change := totalInput - amount
//	// 这里假定找零一定大于交易费，交易费太少的话可能导致交易一直无法确认
//	// 如果change <= fee的话，零钱会转给出块的矿工
//	if change > fee {
//		changePkScript, err := txscript.PayToAddrScript(fromAddr)
//		if err != nil {
//			return nil, errors.Wrap(err, "生成找零地址的脚本失败")
//		}
//		txOut := wire.NewTxOut(change-fee, changePkScript)
//		msgTx.AddTxOut(txOut)
//	}
//
//	// 签署交易
//	// 发送方地址为SegWit的P2WPKH 地址，所以要消费该地址的UTXO，只能通过见证输入进行消费
//	for i, txIn := range msgTx.TxIn {
//		prevOutputScript, err := hex.DecodeString(utxos[i].ScriptPubKey)
//		if err != nil {
//			panic(err)
//		}
//
//		txHash, err := chainhash.NewHashFromStr(utxos[i].TxID)
//		if err != nil {
//			return nil, errors.Wrap(err, "解析交易哈希失败")
//		}
//
//		outPoint := wire.OutPoint{Hash: *txHash, Index: utxos[i].Vout}
//
//		prevOutputFetcher := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
//			outPoint: {
//				Value:    int64(utxos[i].Amount * 1e8),
//				PkScript: prevOutputScript,
//			},
//		})
//
//		sigHashes := txscript.NewTxSigHashes(msgTx, prevOutputFetcher)
//
//		sigScript, err := txscript.WitnessSignature(msgTx, sigHashes, int(utxos[i].Vout), int64(utxos[i].Amount*1e8), prevOutputScript, txscript.SigHashAll, wif.PrivKey, true)
//		if err != nil {
//			return nil, errors.Wrap(err, "签名交易失败")
//		}
//		txIn.Witness = sigScript
//	}
//	return msgTx, nil
//}
