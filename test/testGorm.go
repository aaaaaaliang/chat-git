package main

import (
	"chat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	/*db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})*/
	/*db.AutoMigrate(&models.GroupBasic{})*/
	/*	db.AutoMigrate(&models.Contact{})*/
	db.AutoMigrate(&models.Community{})
	// Create
	/*user := &models.UserBasic{}
	user.Name = "张雪亮"
	db.Create(user)*/
	/*user := &models.UserBasic{
		Name: "张雪亮",
		// 其他字段赋值...
		LoginTime:    time.Now(), // 设置合适的时间值
		LoginOutTime: time.Now(), // 设置合适的时间值
		// 其他字段赋值...
	}
	db.Create(user)*/

	// Read
	/*var product Product*/
	/*db.First(user, 1) // 根据整型主键查找
	/*db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录*/
	/*	fmt.Println(db.First(user, 1))*/

	// Update - 将 product 的 price 更新为 200
	/*db.Model(&product).Update("PassWord", "1234")*/
	// Update - 更新多个字段
	/*db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})*/

	// Delete - 删除 product
	/*db.Delete(&product, 1)*/
}
