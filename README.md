## 功能

### 扫描端口
扫描端口是否能连接，并发邮件通知不能连接的端口

### 判断节点同步
判断节点是否同步，并发邮件通知不能同步的节点

### 查询节点区块高度接口说明
#### 查询指定币种所有节点的区块高度 /getblockheight?cointype=btc
|关键字|类型|描述|
|---|---|---|
|cointype|string|币种类型(BTC,ETH,USDT,ETC,DCR,BTY,EOS,BNB)|
* 返回数据：
```
{
	"id": 0,
	"result": [{
		"cointype": "BTC",
		"error": "",
		"height": 592108,
		"url": "http://192.168.253.105:8802"
	}, {
		"cointype": "BTC",
		"error": "",
		"height": 592108,
		"url": "http://192.168.253.106:8810"
	}],
	"error": null
}
```