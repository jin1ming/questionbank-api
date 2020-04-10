package v1

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
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

	if name == "" || pwd == "" || role == ""{
		c.JSON(200, gin.H{
			"info":"请输入用户名/密码/角色！",
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
	pri, pub, err := RegisterUser(name, role)
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

func DelUser(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("name").(string)
	if userName == ""{
		c.JSON(401, gin.H{
			"info":"您需要先登录!",
		})
		return
	}

	name := c.PostForm("name") //待删除用户
	id, err := GetId(userName, orgName)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"找不到被删除用户!",
		})
		return
	}
	if name == userName || (len(id.Attributes) > 1 && id.Attributes[0].Value == Admin) {
		err = RemoveUser(name, orgName)
		if err != nil {
			c.JSON(200, gin.H{
				"info":"删除用户失败! 请检查您的权利和被删除用户是否存在",
			})
			return
		}
		c.JSON(200, gin.H{
			"info":"删除用户成功!",
		})
		return
	}
}

func PutQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}
	questionId := c.PostForm("id")
	question := c.PostForm("question")
	answer := c.PostForm("answer")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"增加试题失败!code:1",
		})
		return
	}
	var txArgs = [][]byte{[]byte(name), []byte(questionId), []byte(question), []byte(answer)}
	err = executeCC(client,"putQuestion", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"增加试题失败!code:2",
		})
		return
	}
	c.JSON(200, gin.H{
		"info":"增加试题成功!",
	})

}

func DelQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}
	questionId := c.PostForm("id")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"删除试题失败!",
		})
		return
	}

	var txArgs = [][]byte{[]byte(name), []byte(questionId)}
	err = executeCC(client, "delQuestion", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"删除试题失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info":"删除试题成功!",
	})

}

func GetQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}
	questionId := c.PostForm("id")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"获取试题失败!",
		})
		return
	}

	var queryArgs = [][]byte{[]byte(name), []byte(questionId)}
	question, err := queryCC(client, "getQuestion", queryArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "获取试题失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "获取试题成功!",
		"data": json.RawMessage(question),
	})
}

func GetCache(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}

	
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"获取待审核事件失败!",
		})
		return
	}

	events, err := queryCC(client, "getCache", nil)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"获取待审核事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "获取待审核事件成功!",
		"data": json.RawMessage(events),
	})
}

func Approve(c *gin.Context)  {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}
	op := c.PostForm("op")
	questionId := c.PostForm("id")
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"批准事件失败!",
		})
		return
	}
	var queryArgs = [][]byte{[]byte(op), []byte(questionId)}
	err = executeCC(client, "approve", queryArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"批准事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info":"批准事件成功!",
	})
}

func Reject(c *gin.Context)  {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	}else if name == "" {
		c.JSON(200, gin.H{
			"info":"请登录！",
		})
		return
	}
	op := session.Get("op").(string)
	questionId := session.Get("id").(string)
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info":"拒绝事件失败!",
		})
		return
	}
	var queryArgs = [][]byte{[]byte(name),[]byte(op), []byte(questionId)}
	err = executeCC(client, "reject", queryArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info":"拒绝事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info":"拒绝事件成功!",
	})
}
/*
func getLogs(c *gin.Context){
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(Admin), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	response, err := client.Query(

		client.Query()
		{
			ChaincodeID: "qscc",
			Fcn: "invoke",
			Args: integration.ExampleCCQueryArgs("GetChainInfo")
		})
}
 */