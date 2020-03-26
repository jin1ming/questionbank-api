package v1

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var redisdb *redis.Client

var (
	sdk    *fabsdk.FabricSDK
	client *channel.Client
)

func InitRedis(addr string, pwd string, db int) (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       db,  // use default DB
	})

	_, err = redisdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func GetUser(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("user")
	if v == nil {
		c.JSON(204, gin.H{"user": "nil"})
	} else {
		c.JSON(200, gin.H{"user": v.(string)})
	}
}

func AddUser(c *gin.Context) {

}

func UserLogin(c *gin.Context) {

}

func RefreshToken(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}

func UserLogout(c *gin.Context) {

}