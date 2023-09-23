package models

import (
	"chat/utils"
	"fmt"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerID uint
	Img     string
	Desc    string
}

// CreateCommunity 创建群
func CreateCommunity(community Community) (int, string) { //用于响应各种情况
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerID == 0 {
		return -1, "请先登陆"
	}

	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	data := make([]*Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	//utils.DB.Where()
	return data, "查询成功"
}

func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	//contact.TargetId = comId
	contact.Type = 2
	community := Community{}
	//community  根据id和姓名 建立联系
	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	//看contact表是否有群联系
	if community.Name == "" {
		return -1, "没有找到群"
	}
	//判断表字段的时间是不是零时间  来判断是不是已经加入过了
	utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
