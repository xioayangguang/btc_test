package sh

import (
	main2 "btctest/common"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

// P2SH（Pay-to-Script-Hash）
// P2SH 交易允许将资金发送到一个脚本哈希地址（如 3ABC...），实际赎回脚本在花费时提供。
// https://mempool.space/zh/address/3Mg88s8Qjn7Aj9q7wDTRKRnnAbnvFGWkej

func CreateP2SHTransaction() {
	cfg := &chaincfg.TestNet3Params
	// 解码 WIF 私钥
	wif, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	if err != nil {
		log.Fatalf("Failed to decode WIF: %v", err)
	}
	address1, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), cfg)
	pk2, _ := hex.DecodeString("03073d3cf516dceeffaa53a84059fb8701ff5e291b9537457137be851bbc4e5525")
	address2, _ := btcutil.NewAddressPubKey(pk2, cfg)
	pk3, _ := hex.DecodeString("03073d3cf516dceeffaa53a84059fb8701ff5e291b9537457137be851bbc4e5525")
	address3, _ := btcutil.NewAddressPubKey(pk3, cfg)
	script, _ := txscript.MultiSigScript([]*btcutil.AddressPubKey{address1, address2, address3}, 2)
	// 生成 P2SH 地址
	//scriptHash := btcutil.Hash160(script)
	p2shAddr, err := btcutil.NewAddressScriptHash(script, cfg)
	if err != nil {
		log.Fatalf("Failed to create P2SH address: %v", err)
	}
	log.Printf("P2SH testnet address: %s\n", p2shAddr.String())
	// 获取未花费的交易输出（UTXO）
	point, fetcher := main2.GetUnspent(p2shAddr.String())
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
	prevOutput := fetcher.FetchPrevOutput(in.PreviousOutPoint)
	sigScript, err := txscript.SignTxOutput(cfg, tx, 0, prevOutput.PkScript, txscript.SigHashAll, txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return wif.PrivKey, true, nil
	}), nil, nil)

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
