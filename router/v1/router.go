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
		// 登录
		apiv1.POST("/login", apiservice.UserLogin)
		// 登出操作
		apiv1.GET("/logout", apiservice.UserLogout)
		// 刷新token
		apiv1.GET("/token", apiservice.RefreshToken)
		/******************测试API:END******************/
		//注册
		apiv1.POST("/register", apiservice.AddUser)
		// 删除用户
		apiv1.DELETE("/user/:id", apiservice.DeleteUser)
	}

	return r
}