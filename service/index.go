package service

import (
	"chat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /index [get]
func GetIndex(c *gin.Context) {
	/*c.HTML(http.StatusOK, "index.html", nil)*/
	c.File("index.html")

	/*c.JSON(200, gin.H{
		"message": "欢迎光临",
	})*/
}

func ToRegister(c *gin.Context) {
	/*c.HTML(http.StatusOK, "index.html", nil)*/
	/*c.File("index.html")*/
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")

	/*c.JSON(200, gin.H{
		"message": "欢迎光临",
	})*/
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//fmt.Println("ToChat>>>>>>>>", user)
	ind.Execute(c.Writer, user)
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)

}
