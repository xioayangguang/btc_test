包含SH结尾的一般就是多签地址，不是由单个私钥控制

具体的地址类型可以去看https://mempool.space/zh/

P2SH-P2WSH是最复杂的地址，复杂的地址都是嵌套生成的

调试问题对应区块链浏览器排查，查看交易地址类型，查看对应的utxo，查看交易脚本

https://bitcoincore.org/en/segwit_wallet_dev/#creation-of-p2sh-p2wsh-address