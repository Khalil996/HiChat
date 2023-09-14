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
	//好友关系
	relation := v1.Group("relation").Use(middlewear.JWY())
	{
		relation.POST("/list", service.FriendList)
		relation.POST("/add", service.AddFriendByName)
		relation.POST("/new_group", service.NewGroup)
		relation.POST("/group_list", service.GroupList)
		relation.POST("/join_group", service.JoinGroup)
	}

	return router
}
