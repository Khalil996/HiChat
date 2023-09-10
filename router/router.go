package router

import (
	"HiChat/middlewear"
	"HiChat/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	//初始化路由
	router := gin.Default()

	//v1版本
	v1 := router.Group("v1")

	//用户模块
	user := v1.Group("user")
	{
		user.GET("/list", middlewear.JWY(), service.List)
		user.POST("/login_pw", middlewear.JWY(), service.LoginByNameAndPWD)
		user.POST("/create_user", middlewear.JWY(), service.NewUser)
		user.DELETE("/delete_user", middlewear.JWY(), service.DeleteUser)
		user.POST("update_user", middlewear.JWY(), service.UpdateUser)
	}
	return router
}
