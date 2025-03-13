package main

//
//// 从多签钱包取出金额呢
//func main1() {
//	// 使用比特币测试网参数
//	cfg := &chaincfg.TestNet3Params
//
//	// 解码 3 个 WIF 格式的私钥（2-of-3 多签）
//	wif1, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
//	if err != nil {
//		log.Fatalf("Failed to decode WIF 1: %v", err)
//	}
//	wif2, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
//	if err != nil {
//		log.Fatalf("Failed to decode WIF 2: %v", err)
//	}
//	wif3, err := btcutil.DecodeWIF("cViUtGHsa6XUxxk2Qht23NKJvEzQq5mJYQVFRsEbB1PmSHMmBs4T")
//	if err != nil {
//		log.Fatalf("Failed to decode WIF 3: %v", err)
//	}
//
//	// 生成 3 个公钥
//	pubKey1 := wif1.PrivKey.PubKey()
//	pubKey2 := wif2.PrivKey.PubKey()
//	pubKey3 := wif3.PrivKey.PubKey()
//
//	// 创建 2-of-3 多重签名脚本
//	script, err := txscript.MultiSigScript([]*btcutil.AddressPubKey{pubKey1, pubKey2, pubKey3}, 2)
//	if err != nil {
//		log.Fatalf("Failed to create multisig script: %v", err)
//	}
//	// 生成 P2SH 地址
//	scriptHash := btcutil.Hash160(script)
//	p2shAddr, err := btcutil.NewAddressScriptHashFromHash(scriptHash, cfg)
//	if err != nil {
//		log.Fatalf("Failed to create P2SH address: %v", err)
//	}
//	log.Printf("P2SH testnet address: %s\n", p2shAddr.String())
//
//	// 获取未花费的交易输出（UTXO）
//	point, fetcher := GetUnspent(p2shAddr.String())
//
//	// 目标地址
//	destStr := "tb1q4y8u9e0pz7x6w5z3v2c1b0n9m8l7k6j5i4h3g2f1e0d"
//	destAddr, err := btcutil.DecodeAddress(destStr, cfg)
//	if err != nil {
//		log.Fatalf("Failed to decode destination address: %v", err)
//	}
//
//	// 创建交易
//	tx := wire.NewMsgTx(wire.TxVersion)
//
//	// 添加交易输入
//	in := wire.NewTxIn(point, nil, nil)
//	tx.AddTxIn(in)
//
//	// 添加交易输出
//	destScript, err := txscript.PayToAddrScript(destAddr)
//	if err != nil {
//		log.Fatalf("Failed to create destination script: %v", err)
//	}
//	out := wire.NewTxOut(int64(800), destScript) // 800 satoshis
//	tx.AddTxOut(out)
//
//	// 签名交易（需要至少 2 个签名）
//	sigScript, err := SignMultiSigTransaction(tx, script, []*btcutil.WIF{wif1, wif2}, fetcher)
//	if err != nil {
//		log.Fatalf("Failed to sign transaction: %v", err)
//	}
//	tx.TxIn[0].SignatureScript = sigScript
//
//	// 序列化交易
//	var signedTx bytes.Buffer
//	tx.Serialize(&signedTx)
//	finalRawTx := hex.EncodeToString(signedTx.Bytes())
//
//	// 打印最终的签名交易
//	fmt.Printf("Signed Transaction:\n%s\n", finalRawTx)
//}
//
//// SignMultiSigTransaction 对多签交易进行签名
//func SignMultiSigTransaction(tx *wire.MsgTx, script []byte, wifs []*btcutil.WIF, fetcher *txscript.MultiPrevOutFetcher) ([]byte, error) {
//	// 获取交易输入的 UTXO
//	prevOutput := fetcher.FetchPrevOutput(tx.TxIn[0].PreviousOutPoint)
//
//	// 创建签名脚本
//	sigScript, err := txscript.SignTxOutput(
//		&chaincfg.TestNet3Params,
//		tx,
//		0,
//		prevOutput.PkScript,
//		txscript.SigHashAll,
//		txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
//			for _, wif := range wifs {
//				if wif.PrivKey.PubKey().SerializeCompressed() == prevOutput.PkScript[2:22] {
//					return wif.PrivKey, true, nil
//				}
//			}
//			return nil, false, fmt.Errorf("private key not found")
//		}), nil, nil)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return sigScript, nil
//}
