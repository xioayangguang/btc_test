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

//P2WPKH（Pay-to-Witness-PubKey-Hash）
//
//P2WPKH 是 SegWit 地址（如 bc1q... 或 tb1q...），使用隔离见证技术。
//

func CreateP2WPKHTransaction() {
	cfg := &chaincfg.TestNet3Params

	// 解码 WIF 私钥
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}

	// 生成 P2WPKH 地址
	pubKeyHash := btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	p2wpkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2WPKH address: %v", err)
	}
	log.Printf("P2WPKH testnet address: %s\n", p2wpkhAddr.String())

	// 获取未花费的交易输出（UTXO）
	point, fetcher := GetUnspent(p2wpkhAddr.String())

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

	// 签名交易
	prevOutput := fetcher.FetchPrevOutput(in.PreviousOutPoint)
	witness, err := txscript.WitnessSignature(tx, txscript.NewTxSigHashes(tx, fetcher), 0, prevOutput.Value, prevOutput.PkScript, txscript.SigHashAll, wif.PrivKey, true)

	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	tx.TxIn[0].Witness = witness

	// 序列化交易
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())

	// 打印最终的签名交易
	fmt.Printf("Signed Transaction:\n%s\n", finalRawTx)
}
