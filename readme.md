# 基于区块链的题库系统

​	[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)

​	该系统开发目的在：利用智能合约（链码）实现 *RBAC* (基于角色的权限访问控制)，完成对试题库的安全访问，故功能比较简单，可以用来学习fabric开发。

​	系统基于 ``Hyperledger Fabric 1.4.2`` + ``Fabric-sdk-go`` 进行开发。

注意：

- 该系统不包含前端。
- 该系统仅在``linux``(测试了ubuntu18.04、deepin15.11)上可以正常运行，不适用于macOS。

## 环境
​	go1.13+ (GO111MODULE=on)

​	docker latest

​	docker-compose latest

​	mysql5.7
## 初始化配置
- 生成证书以及创世区块

```shell
cd fixtures
./init.sh
```

- 添加hosts

```shell
cd ..
cat hosts >> /etc/hosts
```

- 配置mysql

1. 在mysql中执行``table.sql``
2. 在``controller/api/v1/global.go``中配置用户名密码


## 运行
- 启动区块链及SDK

```shell script
cd fixtures
docker-compose -f docker-compose.yml up -d
cd ..
go run main.go
```
- 关闭区块链

```shell script
docker-compose -f docker-compose.yml down
```

- 若区块链需要重置

**注意：每次重置后必须重启电脑！**

```shell script
cd fixtures
./init.sh
```
