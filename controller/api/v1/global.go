package v1

import (
	"github.com/go-redis/redis"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	redisdb *redis.Client
	sdk    *fabsdk.FabricSDK
	Users map[string]string	// 代替数据库，后期删除使用mysql
)

const (
	Admin = "Admin"
	ccID = "questionbank"
	channelID = "miracle"
	orgName = "Org1"
)