package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"strconv"
)

func GetUser(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("name")
	if v == nil {
		c.JSON(204, gin.H{"user": "nil"})
	} else {
		c.JSON(200, gin.H{"user": v.(string)})
	}
}

func GetAllUsers(c *gin.Context) {

	session := sessions.Default(c)
	nt := session.Get("name")
	if nt == nil  {
		c.JSON(401, gin.H{
			"info": "请登录！",
		})
		return
	}
	// TODO 下面权限检验部分会导致前端偶尔出现401报错,类型stack
	/*
	name := nt.(string)

	id, err := GetId(name, orgName)
	if err != nil || len(id.Attributes) > 1 && id.Attributes[0].Value == Admin {
		c.JSON(401, gin.H{
			"info": "用户不是管理员！",
		})
		return
	}
	*/
	users := getAllUsersFromDb()
	c.JSON(200, gin.H{
		"info": "获取所有用户成功！",
		"data": users,
	})

}

func AddUser(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("pwd")
	role := c.PostForm("role")

	if name == "" || pwd == "" || role == "" {
		c.JSON(200, gin.H{
			"info": "请输入用户名/密码/角色！",
		})
		return
	}
	if checkName(name) {
		c.JSON(200, gin.H{
			"info": "该用户名已被注册！",
		})
		return
	}
	registerDb(name, pwd, role)
	pri, pub, err := RegisterUser(name, role)
	if err != nil {
		panic(err)
		return
	}
	c.JSON(200, gin.H{
		"info": "注册成功！",
		"pri":  pri,
		"pub":  pub,
	})
}

func UserLogin(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("pwd")
	role := c.PostForm("role")
	if name == "" || pwd == "" {
		c.JSON(401, gin.H{
			"info": "请输入用户名密码！",
		})
		return
	}
	if !checkPwd(name, pwd, role) {
		c.JSON(401, gin.H{
			"info": "用户名密码不正确！",
		})
		return
	}
	session := sessions.Default(c)
	session.Set("name", name)
	session.Save()
	c.JSON(200, gin.H{
		"info": "登录成功！",
	})
}

func UserLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("name")
	c.JSON(200, gin.H{
		"info": "用户注销！",
	})
}

func AddPaper(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}
	title := c.PostForm("title")
	owner := name
	qsBytes := []byte(c.PostForm("question_ids"))
	var QuestionIds []string
	err := json.Unmarshal(qsBytes, &QuestionIds)
	if err != nil {
		c.JSON(500, gin.H{
			"info": "Unmarshal wrong!",
			"error": fmt.Sprintf("%s", err),
		})
		return
	}
	addPaper2Db(title, owner, QuestionIds)
	c.JSON(200, gin.H{
		"info":"添加试卷成功！",
	})
}

func DelPaper(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(401, gin.H{
			"info": "请登录！",
		})
		return
	}

	idStr := c.PostForm("id")
	if idStr == "" {
		c.JSON(200, gin.H{
			"info": "未输入试卷id！",
		})
	}

	id, err := strconv.Atoi(idStr)
	CheckErr(err)
	delPaperFromDb(id)
	c.JSON(200, gin.H{
		"info": "删除试卷成功！",
	})
}

func GetAllPapers(c *gin.Context) {
	papers := getAllPapersFromDb()
	c.JSON(200, gin.H{
		"info": "获取所有试卷成功！",
		"data": papers,
	})
}

func GetPaperQuestions(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}

	paperIdStr := c.PostForm("paper_id")
	paperId, err := strconv.ParseInt(paperIdStr, 10, 64)

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "获取试题失败!",
		})
		return
	}

	questionIds := getPaperQuestionsFromDb(paperId)
	//log.Println(questionIds)
	var questions []json.RawMessage
	for _,v := range questionIds {
		//log.Println("v:",v)
		var queryArgs = [][]byte{[]byte(name), []byte(v)}
		question, err := queryCC(client, "getQuestion", queryArgs)
		//log.Println(question)
		CheckErr(err)
		questions = append(questions, json.RawMessage(question))
	}
	questionsBytes, err := json.Marshal(questions)
	c.JSON(200, gin.H{
		"info": "获取试卷所有试题成功！",
		"data": json.RawMessage(questionsBytes),
	})
}

func DelUser(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("name").(string)
	if userName == "" {
		c.JSON(401, gin.H{
			"info": "您需要先登录!",
		})
		return
	}

	name := c.PostForm("name") //待删除用户
	id, err := GetId(userName, orgName)
	if err != nil {
		c.JSON(401, gin.H{
			"info": "找不到该用户!",
		})
		return
	}

	if name == userName || (len(id.Attributes) > 1 && id.Attributes[0].Value == Admin) {
		/*
		err = RemoveUser(name, orgName)
		if err != nil {
			c.JSON(401, gin.H{
				"info": "删除用户失败! 请检查您的权利和被删除用户是否存在",
			})
			return
		}
		 */
		// 从数据库中删除
		delUserFromDb(name)

		c.JSON(200, gin.H{
			"info": "删除用户成功!",
		})
		return
	}


}

func DelPaperItem(c *gin.Context) {
	paperIdStr := c.PostForm("paper_id")
	paperId, err := strconv.ParseInt(paperIdStr, 10, 64)
	var QuestionIds []string
	err = json.Unmarshal([]byte(c.PostForm("question_ids")), &QuestionIds)
	CheckErr(err)
	delPaperItemFromDb(paperId, QuestionIds)
	c.JSON(200, gin.H{
		"info": "移除试卷中部分试题成功！",
	})
}

func PutQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}

	question := c.PostForm("question")
	answer := c.PostForm("answer")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "增加试题失败!code:1",
		})
		return
	}
	var txArgs = [][]byte{[]byte(name),  []byte(question), []byte(answer)}
	err = executeCC(client, "putQuestion", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "增加试题失败!code:2",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "增加试题成功!",
	})

}

func DelQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}
	questionId := c.PostForm("id")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "删除试题失败!",
		})
		return
	}

	var txArgs = [][]byte{[]byte(name), []byte(questionId)}
	err = executeCC(client, "delQuestion", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "删除试题失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "删除试题成功!",
	})

}

func GetQuestion(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}
	questionId := c.PostForm("id")

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "获取试题失败!",
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

func GetAllQuestions(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "获取所有试题失败!",
		})
		return
	}

	var queryArgs = [][]byte{[]byte(name)}
	questions, err := queryCC(client, "getAllQuestions", queryArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "获取所有试题失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "获取所有试题成功!",
		"data": json.RawMessage(questions),
	})
}

func GetCache(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "获取待审核事件失败!",
		})
		return
	}

	events, err := queryCC(client, "getCache", nil)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "获取待审核事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "获取待审核事件成功!",
		"data": json.RawMessage(events),
	})
}

func Approve(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}
	//op := c.PostForm("op")
	questionId := c.PostForm("id")
	op := questionId[6:9]
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "批准事件失败!",
		})
		return
	}
	var txArgs = [][]byte{[]byte(name), []byte(op), []byte(questionId)}
	err = executeCC(client, "approve", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "批准事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "批准事件成功!",
	})
}

func Reject(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}
	//op := c.PostForm("op")
	questionId := c.PostForm("id")
	op := questionId[6:9]
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "拒绝事件失败!",
		})
		return
	}
	var txArgs = [][]byte{[]byte(name), []byte(op), []byte(questionId)}
	err = executeCC(client, "reject", txArgs)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "拒绝事件失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "拒绝事件成功!",
	})
}

func GetLogs(c *gin.Context) {
	session := sessions.Default(c)
	nt := session.Get("name")
	name := c.PostForm("name")
	if nt != nil {
		name = nt.(string)
	} else if name == "" {
		c.JSON(200, gin.H{
			"info": "请登录！",
		})
		return
	}

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)

	if err != nil {
		log.Println("Failed to create new channel client: %s", err)
		c.JSON(200, gin.H{
			"info": "获取日志失败!",
		})
		return
	}

	logs, err := queryCC(client, "getLogs", nil)
	if err != nil {
		c.JSON(200, gin.H{
			"info": "获取日志失败!",
		})
		return
	}
	c.JSON(200, gin.H{
		"info": "获取日志成功!",
		"data": json.RawMessage(logs),
	})
}
