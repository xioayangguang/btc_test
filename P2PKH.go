package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

//P2PKH（Pay-to-PubKey-Hash）
//
//这是最常见的比特币交易类型，使用普通的比特币地址（如 1ABC... 或 tb1q...）。

func main() {
	// 使用比特币测试网参数
	cfg := &chaincfg.TestNet3Params
	// 解码 WIF 格式的私钥
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	// 从私钥生成 P2PKH 地址
	pubKeyHash := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2PKH address: %v", err)
	}
	log.Printf("P2PKH testnet address: %s\n", p2pkhAddr.String())

	// 获取未花费的交易输出（UTXO）
	point, fetcher := GetUnspent(p2pkhAddr.String())

	// 目标地址（接收方地址）
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

	// 签名交易
	prevOutput := fetcher.FetchPrevOutput(in.PreviousOutPoint)
	sigScript, err := txscript.SignatureScript(tx, 0, prevOutput.PkScript, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	tx.TxIn[0].SignatureScript = sigScript

	// 序列化交易
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())

	// 打印最终的签名交易
	fmt.Printf("Signed Transaction:\n%s\n", finalRawTx)
}

// GetUnspent 模拟获取一个未花费的交易输出（UTXO）
func GetUnspent(address string) (*wire.OutPoint, *txscript.MultiPrevOutFetcher) {
	// 模拟一个 UTXO
	txHash, _ := chainhash.NewHashFromStr("7282d54f485561dd21ba22a971b096eb6d0f45ed2fe6bf8c29d87cee162633b4")
	point := wire.NewOutPoint(txHash, uint32(0))

	// 解码地址并生成脚本
	addr, _ := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	script, _ := txscript.PayToAddrScript(addr)

	// 创建一个 UTXO
	output := wire.NewTxOut(int64(1000), script) // 1000 satoshis
	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	fetcher.AddPrevOut(*point, output)

	return point, fetcher
}
