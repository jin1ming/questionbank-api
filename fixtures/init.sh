#!/usr/bin/env bash
#############################
# 作者：靳一鸣                 #
# 时间：2020.3.24             #
#############################
export CHANNEL_NAME="miracle"

echo "证书生成..."
rm -rf crypto-config
rm -rf artifacts/*
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
echo "!!!!!请在之后修改docker-compose.yaml文件!!!!!!!"
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
cryptogen generate --config=./crypto-config.yaml
if [ $? -ne 0 ]; then
    echo "证书生成失败！"
    exit 1
fi

echo "生成创世区块..."
#用于指定启动排序服务
configtxgen -profile MultiNodeEtcdRaft -outputBlock ./artifacts/genesis.block -channelID $CHANNEL_NAME
if [ $? -ne 0 ]; then
    echo "生成创世区块失败！"
    exit 1
fi

echo "生成通道创世区块..."
configtxgen -profile myChannel -outputCreateChannelTx ./artifacts/channel.tx -channelID $CHANNEL_NAME
if [ $? -ne 0 ]; then
    echo "生成通道创世区块失败！"
    exit 1
fi

echo "生成组织锚节点..."
configtxgen -profile myChannel -outputAnchorPeersUpdate ./artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
if [ $? -ne 0 ]; then
    echo "生成组织锚节点失败！"
    exit 1
fi
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
echo "!!!!!请在之后修改docker-compose.yaml文件!!!!!!!"
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"