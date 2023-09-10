package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/middlewear"
	"HiChat/models"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// 用户列表
func List(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

// 密码登录
func LoginByNameAndPWD(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	data, err := dao.FindUserByName(name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "login fail",
		})
		return
	}

	if data.Name == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "user exits",
		})
		return
	}
	ok := common.CheckPassWord(password, data.Salt, data.Password)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "password error",
		})
		return
	}
	Rsp, err := dao.FindUserByNameAndPwd(name, data.Password)
	if err != nil {
		zap.S().Info("error", err)
	}

	//jwt权限认证
	token, err := middlewear.GenerateToken(Rsp.ID, "Khalil")
	if err != nil {
		zap.S().Info("error", err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"token":  token,
		"userID": Rsp.ID,
		"msg":    "login success",
	})

}

// 用户注册
func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("Identity")
	if user.Name == "" || password == "" || repassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "username or password is null",
			"data": user,
		})
		return
	}

	//先判定用户是否存在
	_, err := dao.FindUser(user.Name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "this name is exits",
			"data": user,
		})
		return
	}

	//再判定密码是否一致
	if password != repassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "the password did not match",
			"data": user,
		})
		return
	}
	//生成盐值
	salt := fmt.Sprintf("%d", rand.Int31())

	//加密密码
	user.Password = common.SaltPassWord(password, salt)
	user.Salt = salt
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	dao.CreateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "create user success",
		"data": user,
	})
}

// 跟新用户信息
func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("type fail", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "sign up fail",
		})
		return
	}
	//获取表单内容
	user.ID = uint(id)
	Name := ctx.Request.FormValue("name")
	Password := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("icon")
	gender := ctx.Request.FormValue("gender")
	if Name != "" {
		user.Name = Name
	}
	if Password != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.Password = common.SaltPassWord(Password, salt)
	}
	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}
	//进行参数验证
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("the parameters do not match")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "the parameters do not match",
		})
		return
	}
	Rsp, err := dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("update user fail")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "update user fail",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "update success",
		"data": Rsp.Name,
	})
}

// 账号注销
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("delete user fail ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "delete user fail",
		})
		return
	}
	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("delete user fail", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "delete user fail ",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "delete user success",
	})
}
