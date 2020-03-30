package v1

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var redisdb *redis.Client

var (
	sdk    *fabsdk.FabricSDK
	Users map[string]string
)

func InitRedis(addr string, pwd string, db int) (err error) {

	Users = make(map[string]string)

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
	v := session.Get("name")
	if v == nil {
		c.JSON(204, gin.H{"user": "nil"})
	} else {
		c.JSON(200, gin.H{"user": v.(string)})
	}
}

func AddUser(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("pwd")
	role := c.PostForm("role")

	if name == "" || pwd == "" {
		c.JSON(200, gin.H{
			"info":"请输入用户名密码！",
		})
		return
	}
	if Users[name] != "" {
		c.JSON(200,gin.H{
			"info":"该用户名已被注册！",
		})
		return
	}
	Users[name] = pwd
	pri, pub, err := RegisterUser(name, "Org1", role)
	if err != nil {
		panic(err)
		return
	}
	c.JSON(200, gin.H{
		"info": "注册成功！",
		"pri": pri,
		"pub": pub,
	})
}

func UserLogin(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("pwd")
	if name == "" || pwd == "" {
		c.JSON(200, gin.H{
			"info":"请输入用户名密码！",
		})
		return
	}
	if Users[name] != pwd {
		c.JSON(200, gin.H{
			"info":"用户名密码不正确！",
		})
		return
	}
	session := sessions.Default(c)
	session.Set("name", name)
	session.Save()
	c.JSON(200, gin.H{
		"info":"登录成功！",
	})
}

func RefreshToken(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}

func UserLogout(c *gin.Context) {

}