package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
	"go.uber.org/zap"
)

// 好友列表
func GetFriendList(userId uint) (*[]models.UserBasic, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id=? and type=1", userId).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("dont find relation")
		return nil, errors.New("dont find relationship")
	}

	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetID)
	}
	user := make([]models.UserBasic, 0)
	if tx := global.DB.Where("id in ?", userID).Find(&user); tx.RowsAffected == 0 {
		zap.S().Info("dont find relationship")
		return nil, errors.New("dont find friend")
	}
	return &user, nil
}

// 通过id添加好友，数据库进行俩次添加（数据库事务特征）一致性
func AddFriend(userID, TargetId uint) (int, error) {

	//判定userid和targetid是否一致
	if userID == TargetId {
		return -2, errors.New("userID and TargetId equal")
	}

	//通过id查询用户
	targetUser, err := FindUserID(TargetId)
	if err != nil {
		return -1, errors.New("dont find user")
	}
	if targetUser.ID == 0 {
		zap.S().Info("dont find user")
		return -1, errors.New("dont find user")
	}

	relation := models.Relation{}
	if tx := global.DB.Where("owner_id =? and target_id=? and type=1", userID, TargetId).Find(&relation); tx.RowsAffected == 1 {
		zap.S().Info("this friend exists")
		return 0, errors.New("this friend exists")
	}
	if tx := global.DB.Where("owner_id =? and target_id=? and type=1", TargetId, userID).Find(&relation); tx.RowsAffected == 1 {
		zap.S().Info("this friend exists")
		return 0, errors.New("this friend exists")
	}

	//开始事务
	tx := global.DB.Begin()

	relation.OwnerId = userID
	relation.TargetID = targetUser.ID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("create fail")

		//事务回滚
		tx.Rollback()
		return -1, errors.New("create relation fail")
	}

	//提交事务
	tx.Commit()

	return 1, nil
}

// 通过昵称添加（昵称获取id，用id添加用户）
func AddFriendByName(userId uint, targetName string) (int, error) {
	user, err := FindUser(targetName)
	if err != nil {
		return -1, errors.New("this user no exists")
	}
	if user.ID == 0 {
		zap.S().Info("dont find user")
		return -1, errors.New("this user not exists")
	}
	return AddFriend(userId, user.ID)
}
