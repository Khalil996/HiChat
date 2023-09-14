package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type user struct {
	Name     string
	Avatar   string
	Gender   string
	Phone    string
	Email    string
	Identity string
}

// 好友列表
func FriendList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Request.FormValue("userId"))
	users, err := dao.GetFriendList(uint(id))
	if err != nil {
		zap.S().Info("get friend list fail")
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "friend is nil",
		})
		return
	}
	infos := make([]user, 0)
	for _, v := range *users {
		info := user{
			Name:     v.Name,
			Avatar:   v.Avatar,
			Gender:   v.Gender,
			Phone:    v.Phone,
			Email:    v.Email,
			Identity: v.Identity,
		}
		infos = append(infos, info)
	}
	common.RespOKList(ctx.Writer, infos, len(infos))
}

// 添加好友
func AddFriendByName(ctx *gin.Context) {
	user := ctx.PostForm("userId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("type fail", err)
		return
	}
	tar := ctx.PostForm("targetName")
	target, err := strconv.Atoi(tar)
	if err != nil {
		code, err := dao.AddFriendByName(uint(userId), tar)
		if err != nil {
			HandleErr(code, ctx, err)
			return
		}
	} else {
		code, err := dao.AddFriend(uint(userId), uint(target))
		if err != nil {
			HandleErr(code, ctx, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "add friend success",
	})

}
func HandleErr(code int, ctx *gin.Context, err error) {
	switch code {
	case -1:
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1, //  0成功   -1失败
			"msg":  err.Error(),
		})
	case 0:
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1, //  0成功   -1失败
			"msg":  "this user exists",
		})
	case -2:
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1, //  0成功   -1失败
			"msg":  "dont add yourself",
		})

	}
}

// 新建群
func NewGroup(ctx *gin.Context) {
	owner := ctx.PostForm("ownerId")
	ownerId, err := strconv.Atoi(owner)
	if err != nil {
		zap.S().Info("owner type fail", err)
		return
	}
	ty := ctx.PostForm("cate")
	Type, err := strconv.Atoi(ty)
	if err != nil {
		zap.S().Info("ty type fail", err)
		return
	}
	img := ctx.PostForm("icon")
	name := ctx.PostForm("name")
	desc := ctx.PostForm("desc")

	community := models.Community{}
	if ownerId == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "please login",
		})
		return
	}
	if name == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "name dont nil",
		})
		return
	}
	if img != "" {
		community.Image = img
	}
	if desc != "" {
		community.Desc = desc
	}
	community.Name = name
	community.Type = Type
	community.ID = uint(ownerId)

	code, err := dao.CreateCommunity(community)
	if err != nil {
		HandleErr(code, ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "create success",
	})
}

// 群列表
func GroupList(ctx *gin.Context) {
	owner := ctx.PostForm("ownerId")
	ownerId, err := strconv.Atoi(owner)
	if err != nil {
		zap.S().Info("owner type fail", err)
		return
	}
	if ownerId == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "please login",
		})
		return
	}
	rsp, err := dao.GetCommunityList(uint(ownerId))
	if err != nil {
		zap.S().Info("get group fail", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "yon dont join any group",
		})
		return
	}
	common.RespOKList(ctx.Writer, rsp, len(*rsp))

}

// 加入群聊
func JoinGroup(ctx *gin.Context) {
	comInfo := ctx.PostForm("comId")
	if comInfo == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "community name dont nil",
		})
		return
	}
	user := ctx.PostForm("userId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("userid type fail", err)
		return
	}
	code, err := dao.JoinCommunity(uint(userId), comInfo)
	if err != nil {
		HandleErr(code, ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "join success",
	})
}
