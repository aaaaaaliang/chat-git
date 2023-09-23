package models

import (
	"chat/utils"
	"fmt"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //对应的类型1好友 2群主 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userID int) []UserBasic {
	contacts := make([]Contact, 0)
	//存储目标id的用户
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=1", userID).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>>>>>", v)
		// 将目标用户的 ID 添加到 objIds 切片中
		objIds = append(objIds, uint64(v.TargetId))
	}
	// 创建一个空的 UserBasic 切片来存储查询结果
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友
func AddFriend(userId uint, targetName string) (int, string) {
	/*user := UserBasic{}*/
	if targetName != "" {
		//找到对应的人
		targetUser := FindByName(targetName)
		if targetUser.Salt != "" {
			if targetUser.ID == userId {
				return -1, "不能自己添加自己"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id = ? and target_id =? and type=1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "该用户重复添加"
			}

			//开启事物   要么全部执行成功，要么全部不执行
			tx := utils.DB.Begin()
			//事物一旦开始 无论出现什么异常 最终都会rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetUser.ID
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}

			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "找不到此用户"
	}
	return -1, "好友ID不能为空"
}
func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
