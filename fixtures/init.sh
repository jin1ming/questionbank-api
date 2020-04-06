#!/usr/bin/env bash
#############################
# 文件名：init.sh             #
# 作者：dakini_wind           #
# 时间：2020.3.25             #
#############################
export CHANNEL_NAME="miracle" #通道名

echo "证书生成..."
rm -rf crypto-config  #创世区块、通道区块目录
rm -rf artifacts/*    #存放证书的目录
cryptogen generate --config=./crypto-config.yaml
if [ $? -ne 0 ]; then
    echo "证书生成失败！"
    exit 1
fi

echo "自动修改docker-compose.yaml文件..."
pushd .
cd crypto-config/peerOrganizations/org1.questionbank.com/ca/
filename=`ls *_sk`
echo "获取到ca-key："$filename
popd
pre="FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/"
sed -i "s#$pre.*#$pre$filename#" docker-compose.yml
if [ $? -ne 0 ]; then
    echo "自动修改出错！code:1"
    exit 1
fi
pre="FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/"
sed -i "s#$pre.*#$pre$filename#" docker-compose.yml
if [ $? -ne 0 ]; then
    echo "自动修改出错！code:2"
    exit 1
fi


pre="--ca.keyfile /etc/hyperledger/fabric-ca-server-config/"
sed -i "s#$pre.*_sk#$pre$filename#" docker-compose.yml
if [ $? -ne 0 ]; then
    echo "自动修改出错！code:3"
    exit 1
fi

echo "生成创世区块..."
#用于指定启动排序服务
configtxgen -profile MultiNodeEtcdRaft -outputBlock ./artifacts/genesis.block
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
echo "!!!!!!!自动化生成区块、证书以及修改成功!!!!!!!!!!"
echo "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"