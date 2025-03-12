package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

// https://decision01.com/post/c19159ae
func main() {
	cfg := &chaincfg.TestNet3Params
	wif, _ := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
	taprootAddr, _ := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wif.PrivKey.PubKey())),
		&chaincfg.TestNet3Params,
	)
	log.Printf("Taproot testnet address: %s\n", taprootAddr.String())

	point, fetcher := GetUnspent(taprootAddr.String())

	destStr := "tb1pvwak065fek4y0mup9p4l7t03ey2nu8as7zgcrlgm9mdfl8gs5rzss490qd"
	byteAddr, _ := DecodeTaprootAddress(destStr, cfg)

	tx := wire.NewMsgTx(wire.TxVersion)

	in := wire.NewTxIn(point, nil, nil)
	tx.AddTxIn(in)

	out := wire.NewTxOut(int64(800), byteAddr)
	tx.AddTxOut(out)

	prevOutput := fetcher.FetchPrevOutput(in.PreviousOutPoint)
	witness, _ := txscript.TaprootWitnessSignature(
		tx,
		txscript.NewTxSigHashes(tx, fetcher),
		0,
		prevOutput.Value,
		prevOutput.PkScript,
		txscript.SigHashDefault,
		wif.PrivKey,
	)
	tx.TxIn[0].Witness = witness

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Printf("Signed Transaction:\n%s", finalRawTx)
}

//
//func GetUnspent(address string) (*wire.OutPoint, *txscript.MultiPrevOutFetcher) {
//	txHash, _ := chainhash.NewHashFromStr(
//		"7282d54f485561dd21ba22a971b096eb6d0f45ed2fe6bf8c29d87cee162633b4")
//	point := wire.NewOutPoint(txHash, uint32(0))
//	script, _ := hex.DecodeString("51208b63f2ee8d7a385e12c0e0f7599cd86ef6e2aed7b9e033762afb177f16c2f309")
//	output := wire.NewTxOut(int64(1000), script)
//	fetcher := txscript.NewMultiPrevOutFetcher(nil)
//	fetcher.AddPrevOut(*point, output)
//	return point, fetcher
//}

func DecodeTaprootAddress(strAddr string, cfg *chaincfg.Params) ([]byte, error) {
	taprootAddr, err := btcutil.DecodeAddress(strAddr, cfg)
	if err != nil {
		return nil, err
	}
	byteAddr, err := txscript.PayToAddrScript(taprootAddr)
	if err != nil {
		return nil, err
	}
	return byteAddr, nil
}
