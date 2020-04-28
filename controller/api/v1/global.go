package v1

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	//redisdb *redis.Client
	sdk    *fabsdk.FabricSDK
	//Users map[string]string	// 代替数据库，后期删除使用mysql
)

const (
	// Fabric 配置
	Admin = "Admin"
	ccID = "questionbank"
	channelID = "miracle"
	orgName = "org1"
	// mysql 配置
	DBHostsIp  = "127.0.0.1:3306"
	DBUserName = "root"
	DBPassWord = "123456"
	DBName     = "questionbank"
)