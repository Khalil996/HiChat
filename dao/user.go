package dao

import (
	"HiChat/common"
	"HiChat/global"
	"HiChat/models"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// 用户列表
func GetUserList() ([]*models.UserBasic, error) {
	var list []*models.UserBasic
	if tx := global.DB.Find(&list); tx.RowsAffected == 0 {
		return nil, errors.New("get UserList Fail")
	}
	return list, nil
}

// 用户名密码查询
func FindUserByNameAndPwd(name, password string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=? and pass_word=?", name, password).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("find fail")
	}
	//登录识别
	t := strconv.Itoa(int(time.Now().Unix()))

	temp := common.Md5encoder(t)

	if tx := global.DB.Model(&user).Where("id=?", user.ID).Update("identity", temp); tx.RowsAffected == 0 {
		return nil, errors.New("identity fail")
	}
	return &user, nil
}

// 登录时查找用户
func FindUserByName(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=?", name).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("find username fail")
	}
	return &user, nil
}

// 注册时查找用户
func FindUser(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=?", name).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("the user exists")
	}
	return &user, nil
}

// ID查询用户
func FindUserID(ID uint) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where(ID).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("not user")
	}
	return &user, nil

}

// 根据电话查询用户
func FindUserByPhone(phone string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("phone=?", phone).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("not user")
	}
	return &user, nil
}

// 根据邮件查询用户
func FindUserByEmail(email string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("email=?", email).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("not find")
	}
	return &user, nil
}

// 新建用户
func CreateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Create(&user)
	if tx.RowsAffected == 0 {
		zap.S().Info("create fail")
		return nil, errors.New("create fail")
	}
	return &user, nil
}

// 删除用户
func DeleteUser(user models.UserBasic) error {
	tx := global.DB.Delete(&user)
	if tx.RowsAffected == 0 {
		zap.S().Info("delete fail")
		return errors.New("delete fail")
	}
	return nil
}
func UpdateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Salt:     user.Salt,
	})
	if tx.RowsAffected == 0 {
		zap.S().Info("Update fail")
		return nil, errors.New("update fail")
	}
	return &user, nil
}
