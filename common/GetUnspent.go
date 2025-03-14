package common

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

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
