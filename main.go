package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/github-user/IMv_3/controllers"
	"github.com/github-user/IMv_3/kafka"
	"github.com/github-user/IMv_3/models"
)

func init() {
	models.Setup()
	kafka.Setup()
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionid", store))
	api := r.Group("/api")
	{
		usr := api.Group("/user")
		{
			//注册账号
			usr.POST("/register", controllers.Register)
			//登录账号,返回sessionid
			usr.POST("/login", controllers.Login)

		}
		chat := api.Group("/chat")
		{
			//建立websocket连接
			chat.GET("/connect", controllers.ConnectWeb)
			//发消息
			chat.POST("/send_message", controllers.SendMessage)
			//查询会话列表
			chat.GET("", controllers.FindChats)
			//修改会话状态
			chat.PUT("/:id", controllers.UpdateChat)
			//查询历史聊天记录
			chat.GET("/find_messages", controllers.FindMessages)
			chat.Use(controllers.Authorize())
		}

	}
	r.Run(":9999")
}
