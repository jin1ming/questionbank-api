# 基于区块链的题库系统
## 1.环境
go1.13+ (GO111MODULE=on)

docker lasest

docker-compose lasest
## 2.运行步骤
启动区块链及SDK：
```shell script
cd fixtures
docker-compose -f docker-compose.yml up -d
cd ..
go run main.go
```
关闭区块链：
```shell script
docker-compose -f docker-compose.yml up -down
```

若区块链需要重置：
```shell script
cd fixtures
./init.sh
```