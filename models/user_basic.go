package models

import (
	"chat/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model    //gorm包里的属性
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Duration
	LoginOutTime  time.Time `gorm:"column:login_out_time;default:null" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	//返回我想要的表名
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:          user.Name,
		PassWord:      user.PassWord,
		Email:         user.Email,
		Phone:         user.Phone,
		LoginOutTime:  user.LoginOutTime,  // 更新下线时间
		HeartbeatTime: user.HeartbeatTime, // 更新在线时长
		ClientIp:      user.ClientIp,
		IsLogout:      user.IsLogout,
	})
}

// 登陆
func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Where("name= ? and pass_word=?", name, password).First(&user)
	utils.DB.Model(&user).Where("id=?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name= ?", name).First(&user)
	return user
}
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Phone= ?", phone).First(&user)
	return user
}
func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Email= ?", email).First(&user)
	return user
}

//查找某个用户

/*
	func FindByID(id uint) UserBasic {
		user := UserBasic{}
		utils.DB.Where("id =?", id).First(&user)
		return user
	}
*/
func FindByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name =?", name).First(&user)
	return user
}

// 查找某个用户
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
