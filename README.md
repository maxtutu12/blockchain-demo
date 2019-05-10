# 用golang 构建一个简单的区块链

最近在学习区块链技术的一些知识，通过用golang来实践以加深理解

项目参考了 [用Python从零开始创建区块链](https://learnblockchain.cn/2017/10/27/build_blockchain_by_python/)

## 运行项目

>拉取项目
```
git clone https://github.com/xprint120/blockchain-demo.git
```
>进入目录, 运行节点
```
go run cmd/main.go -port=5000
go run cmd/main.go -port=5001
```

可以用[postmain](https://www.getpostman.com/downloads/)工具来进行测试

挖矿通过发送get请求:
```
http://localhost:5000/mine
```
添加新的交易发送post请求:
```
http://localhost:5000/transactions/new

{
 "sender": "my address",
 "recipient": "someone else's address",
 "amount": 5
}
```
获取块信息发送get请求:
```
http://localhost:5000/chain
```
在5001节点上执行注册节点,发送post请求:
```
http://localhost:5001/nodes/register

{
	"node": ["http://localhost:5000"]
}
```
在节点5001上达成共识,通过发送get请求：
```
http://localhost:5001/nodes/resolve
```