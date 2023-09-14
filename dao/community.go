package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
)

// 新建群
func CreateCommunity(community models.Community) (int, error) {
	com := models.Community{}
	//查询群是否存在，考虑群名字可以重复所以这里通过群id来查找
	if tx := global.DB.Where("ID=?", community.ID).First(com); tx.RowsAffected == 1 {
		return -1, errors.New("this community exists")
	}
	tx := global.DB.Begin()
	if t := tx.Create(&community); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("create fail")
	}
	relation := models.Relation{}
	relation.OwnerId = community.OwnerId
	relation.ID = community.ID
	relation.Type = 2
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("create fail")
	}

	tx.Commit()
	return 0, nil
}

// 获取群列表
func GetCommunityList(ownerId uint) (*[]models.Community, error) {

	//获取群id
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id=? and type=2", ownerId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("community is not exists")
	}
	communityID := make([]uint, 0)

	for _, v := range relation {
		cid := v.TargetID
		communityID = append(communityID, cid)
	}
	community := make([]models.Community, 0)
	if tx := global.DB.Where("id in ?", communityID).Find(&community); tx.RowsAffected == 0 {
		return nil, errors.New("get community fail")
	}
	return &community, nil
}

// 添加群聊
func JoinCommunity(ownerId uint, cname string) (int, error) {
	community := models.Community{}
	//查询群是否存在
	if tx := global.DB.Where("name =?", cname).First(&community); tx.RowsAffected == 0 {
		return -1, errors.New("community is not exists")
	}

	//重复加群
	relation := models.Relation{}
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type =2 ", ownerId, community.ID).First(&relation); tx.RowsAffected == 1 {
		return -1, errors.New("this community is exists")
	}

	relation = models.Relation{}
	relation.OwnerId = ownerId
	relation.TargetID = community.ID
	relation.Type = 2

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		return -1, errors.New("join fail")
	}
	return 0, nil
}
