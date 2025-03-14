package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestAAA(t *testing.T) {
	utxos, err := GetUTXOs("mkWvgCpdMVq6xECWSKjRv5VaPhzUidyZEE")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(err)
	a, _ := json.Marshal(utxos)
	fmt.Println(string(a))

	utxos, err = GetUTXOs("mvSQRoib2ge45xNu3UWLrEfMYZWWn1oBa1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(err)
	a, _ = json.Marshal(utxos)
	fmt.Println(string(a))
}

func TestCC(t *testing.T) {
	txOut, pkScript, err := BuildTxOut("mkWvgCpdMVq6xECWSKjRv5VaPhzUidyZEE", 1000, chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pkScript))
	privKeyWif, err := btcutil.DecodeWIF("92exWYZpKJQScvhuYqszQFMfahJzQCwZCBH5N6R6wRhSPok9rti")
	if err != nil {
		panic(err)
	}
	msgTx, err := BuildTxIn(privKeyWif, 1000, txOut, &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	txHash := msgTx.TxHash()
	fmt.Println(txHash)
	a, _ := txHash.MarshalJSON()
	fmt.Println(string(a))
	buf := &bytes.Buffer{}
	err = msgTx.Serialize(buf)
	if err != nil {
		panic(err)
	}
	txHex := fmt.Sprintf("%x", hex.EncodeToString(buf.Bytes()))
	fmt.Println(txHex)

	txHash1, err := client.SendRawTransaction(msgTx, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash1)
	//cb95082e16fb8ee47d73f847a2b442f69510a7af619b8074e265805762b00c2c  //成功

	//txHex := "0100000001c997a5e56e104102fa209c6a852dd90660a20b2d9c352423edce25857fcd3704000000004847304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901ffffffff0200ca9a3b00000000434104ae1a62fe09c5f51b13905f07f06b99a2f7159b2225f374cd378d71302fa28414e7aab37397f554a7df5f142c21c1b7303b8a0626f1baded5c72a704f7e6cd84cac00286bee0000000043410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac00000000"
	//// 将十六进制字符串解码为字节
	//txBytes, err := hex.DecodeString(txHex)
	//if err != nil {
	//	fmt.Println("Error decoding hex:", err)
	//	return
	//}
	//
	//// 解析交易
	//var tx wire.TxIn
	//err = tx.Deserialize(bytes.NewReader(txBytes))
	//if err != nil {
	//	fmt.Println("Error deserializing transaction:", err)
	//	return
	//}
}

func TestBB(t *testing.T) {
	wifKey, address, _ := GenerateBTCTest() // 测试地址
	// wifKey, address, _ := GenerateBTC() // 正式地址
	fmt.Println(address, wifKey)
	//mkWvgCpdMVq6xECWSKjRv5VaPhzUidyZEE
	//92NEVHkAkJ2s3aEDJtFrzJKJig5omLvhabkGkxUVu4GJv5ptvMr

	// mvSQRoib2ge45xNu3UWLrEfMYZWWn1oBa1
	// 92exWYZpKJQScvhuYqszQFMfahJzQCwZCBH5N6R6wRhSPok9rti
}
