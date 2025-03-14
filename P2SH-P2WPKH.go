package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

// 3开始的私人地址，(也就是兼容隔离见证地址)
func main() {
	// 使用比特币测试网参数
	cfg := &chaincfg.TestNet3Params

	// 解码 WIF 格式的私钥
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}

	// 生成 P2WPKH 脚本
	pubKeyHash := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	p2wpkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2WPKH address: %v", err)
	}
	p2wpkhScript, err := txscript.PayToAddrScript(p2wpkhAddr)
	if err != nil {
		log.Fatalf("Failed to create P2WPKH script: %v", err)
	}

	// 生成 P2SH-P2WPKH 地址
	p2shAddr, err := btcutil.NewAddressScriptHash(p2wpkhScript, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2SH address: %v", err)
	}
	log.Printf("P2SH-P2WPKH testnet address: %s\n", p2shAddr.String())

	// 获取未花费的交易输出（UTXO）
	point, fetcher := GetUnspent(p2shAddr.String())

	// 创建交易
	tx := wire.NewMsgTx(wire.TxVersion)
	// 添加交易输入
	in := wire.NewTxIn(point, nil, nil)
	tx.AddTxIn(in)

	// 目标地址
	destStr := "tb1q4y8u9e0pz7x6w5z3v2c1b0n9m8l7k6j5i4h3g2f1e0d"
	destAddr, err := btcutil.DecodeAddress(destStr, cfg)
	if err != nil {
		log.Fatalf("Failed to decode destination address: %v", err)
	}

	// 添加交易输出
	destScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		log.Fatalf("Failed to create destination script: %v", err)
	}
	out := wire.NewTxOut(int64(800), destScript) // 800 satoshis
	tx.AddTxOut(out)

	// 签名交易
	witness, sigScript, err := SignP2SHP2WPKHTransaction(tx, p2wpkhScript, wif, fetcher)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	tx.TxIn[0].SignatureScript = sigScript //todo 这里不知道正确与否，需要给SignatureScript赋值与否
	tx.TxIn[0].Witness = witness

	// 序列化交易
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())
	// 打印最终的签名交易
	fmt.Printf("Signed Transaction:\n%s\n", finalRawTx)
}

// SignP2SHP2WPKHTransaction 对 P2SH-P2WPKH 交易进行签名
func SignP2SHP2WPKHTransaction(tx *wire.MsgTx, p2wpkhScript []byte, wif *btcutil.WIF, fetcher *txscript.MultiPrevOutFetcher) (wire.TxWitness, []byte, error) {
	// 获取交易输入的 UTXO
	prevOutput := fetcher.FetchPrevOutput(tx.TxIn[0].PreviousOutPoint)

	// 创建签名哈希
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)

	// 生成见证签名
	witnessSig, err := txscript.WitnessSignature(
		tx, sigHashes, 0, prevOutput.Value, p2wpkhScript,
		txscript.SigHashAll, wif.PrivKey, true,
	)
	if err != nil {
		return nil, nil, err
	}

	// 生成签名脚本（包含 P2WPKH 赎回脚本）
	sigScript, err := txscript.NewScriptBuilder().AddData(p2wpkhScript).Script()
	if err != nil {
		return nil, nil, err
	}

	return witnessSig, sigScript, nil
}
