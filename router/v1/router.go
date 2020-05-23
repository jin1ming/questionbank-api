package v1

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	apiservice "questionbank-api/controller/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 调试需要session，之后移除
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	gin.SetMode(gin.DebugMode)

	apiv1 := r.Group("/api/v1")
	{
		/*****************测试API:Start*****************/
		//获取用户
		apiv1.GET("/user", apiservice.GetUser)

		// 获取所有用户
		apiv1.GET("/get_all_users", apiservice.GetAllUsers)
		// 登出操作
		apiv1.GET("/logout", apiservice.UserLogout)
		// 登录
		apiv1.POST("/login", apiservice.UserLogin)
		// 注册
		apiv1.POST("/register", apiservice.AddUser)
		// 删除用户
		apiv1.POST("/delete", apiservice.DelUser)
		// 增加试题
		apiv1.POST("/put_question", apiservice.PutQuestion)
		// 删除试题
		apiv1.POST("/del_question", apiservice.DelQuestion)
		// 查询试题
		apiv1.POST("/get_question", apiservice.GetQuestion)
		// 获取所有试题
		apiv1.POST("/get_all_questions", apiservice.GetAllQuestions)
		// 获取所有试卷
		apiv1.GET("/get_all_papers", apiservice.GetAllPapers)
		// 添加试卷
		apiv1.POST("/add_paper", apiservice.AddPaper)
		// 获取某试卷所有试题
		apiv1.POST("/get_paper", apiservice.GetPaperQuestions)
		// 删除某试卷中的部分试题(暂不实现)
		apiv1.POST("/del_paper_item", apiservice.DelPaperItem)
		// 获取待审核事件
		apiv1.POST("/get_cache", apiservice.GetCache)
		// 批准事件
		apiv1.POST("/approve", apiservice.Approve)
		// 拒绝事件
		apiv1.POST("/reject", apiservice.Reject)
		// 获取日志
		apiv1.POST("/get_logs", apiservice.GetLogs)
	}

	return r
}