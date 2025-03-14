package main

import (
	main2 "btctest"
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

//
//P2WSH（Pay-to-Witness-Script-Hash）
//
//P2WSH 是 SegWit 的多重签名地址（如 bc1q... 或 tb1q...），使用隔离见证技术
//
//

func CreateP2WSHTransaction() {
	cfg := &chaincfg.TestNet3Params

	// 解码 WIF 私钥
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}

	// 创建多重签名脚本（2-of-3）
	//pubKey1 := wif.PrivKey.PubKey()
	//pubKey2, _ := btcec.NewPrivateKey()
	//pubKey3, _ := btcec.NewPrivateKey()

	address1, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), cfg)
	pk2, _ := hex.DecodeString("03073d3cf516dceeffaa53a84059fb8701ff5e291b9537457137be851bbc4e5525")
	address2, _ := btcutil.NewAddressPubKey(pk2, cfg)
	pk3, _ := hex.DecodeString("03073d3cf516dceeffaa53a84059fb8701ff5e291b9537457137be851bbc4e5525")
	address3, _ := btcutil.NewAddressPubKey(pk3, cfg)

	script, _ := txscript.MultiSigScript([]*btcutil.AddressPubKey{address1, address2, address3}, 2)

	// 生成 P2WSH 地址
	scriptHash := sha256.Sum256(script)
	p2wshAddr, err := btcutil.NewAddressWitnessScriptHash(scriptHash[:], cfg)
	if err != nil {
		log.Fatalf("Failed to create P2WSH address: %v", err)
	}
	log.Printf("P2WSH testnet address: %s\n", p2wshAddr.String())

	// 获取未花费的交易输出（UTXO）
	point, fetcher := main2.GetUnspent(p2wshAddr.String())

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
