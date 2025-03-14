package sh

import (
	"btctest/common"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

// https://mempool.space/zh/address/3MaE1ACBvr5SsQFHkpVU12tMDkvZ9Yxruu
func main() {
	// 使用比特币测试网参数
	cfg := &chaincfg.TestNet3Params
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	address1, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), cfg)

	wif1, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	address2, _ := btcutil.NewAddressPubKey(wif1.PrivKey.PubKey().SerializeUncompressed(), cfg)

	wif2, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	address3, _ := btcutil.NewAddressPubKey(wif2.PrivKey.PubKey().SerializeUncompressed(), cfg)

	script, _ := txscript.MultiSigScript([]*btcutil.AddressPubKey{address1, address2, address3}, 2)

	// 生成 P2WSH 脚本
	witnessScriptHash := sha256.Sum256(script)
	p2wshAddr, err := btcutil.NewAddressWitnessScriptHash(witnessScriptHash[:], cfg)
	if err != nil {
		log.Fatalf("Failed to create P2WSH script: %v", err)
	}
	p2wshScript, err := txscript.PayToAddrScript(p2wshAddr)
	if err != nil {
		log.Fatalf("Failed to create P2WSH script: %v", err)
	}
	// 生成 P2SH-P2WSH 地址
	p2shAddr, err := btcutil.NewAddressScriptHash(p2wshScript, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2SH address: %v", err)
	}
	log.Printf("P2SH-P2WSH testnet address: %s\n", p2shAddr.String())

	// 获取未花费的交易输出（UTXO）
	point, fetcher := common.GetUnspent(p2shAddr.String())

	// 目标地址
	destStr := "tb1q4y8u9e0pz7x6w5z3v2c1b0n9m8l7k6j5i4h3g2f1e0d"
	destAddr, err := btcutil.DecodeAddress(destStr, cfg)
	if err != nil {
		log.Fatalf("Failed to decode destination address: %v", err)
	}

	// 创建交易
	tx := wire.NewMsgTx(wire.TxVersion)

	// 添加交易输入
	in := wire.NewTxIn(point, nil, nil)
	tx.AddTxIn(in)

	// 添加交易输出
	destScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		log.Fatalf("Failed to create destination script: %v", err)
	}
	out := wire.NewTxOut(int64(800), destScript) // 800 satoshis
	tx.AddTxOut(out)

	// 签名交易（需要至少 2 个签名）
	witness, sigScript, err := SignP2SHP2WSHTransaction(tx, script, []*btcutil.WIF{wif1, wif2}, fetcher)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	tx.TxIn[0].SignatureScript = sigScript
	tx.TxIn[0].Witness = witness

	// 序列化交易
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())

	// 打印最终的签名交易
	fmt.Printf("Signed Transaction:\n%s\n", finalRawTx)
}

// SignP2SHP2WSHTransaction 对 P2SH-P2WSH 交易进行签名
func SignP2SHP2WSHTransaction(tx *wire.MsgTx, script []byte, wifs []*btcutil.WIF, fetcher *txscript.MultiPrevOutFetcher) (wire.TxWitness, []byte, error) {
	// 获取交易输入的 UTXO
	prevOutput := fetcher.FetchPrevOutput(tx.TxIn[0].PreviousOutPoint)

	// 创建签名哈希
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)

	// 生成见证签名
	witnessSig, err := txscript.WitnessSignature(tx, sigHashes, 0, prevOutput.Value, script, txscript.SigHashAll, wifs[0].PrivKey, true)
	if err != nil {
		return nil, nil, err
	}

	// 添加第二个签名
	witnessSig2, err := txscript.WitnessSignature(tx, sigHashes, 0, prevOutput.Value, script, txscript.SigHashAll, wifs[1].PrivKey, true)
	if err != nil {
		return nil, nil, err
	}

	_ = witnessSig2

	//todo 签名是肯定不对的

	// 生成签名脚本（包含 P2WSH 赎回脚本）
	sigScript, err := txscript.NewScriptBuilder().AddData(script).Script()
	if err != nil {
		return nil, nil, err
	}

	return witnessSig, sigScript, nil
}
