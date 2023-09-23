package service

import (
	"chat/models"
	"chat/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @param email query string false "邮箱"
// @param phone query string false "电话"
// @param identity query string false "身份证"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	/*user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")*/

	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	//生成六位数的随机密码
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"message": "用户名或密码不能为空!",
		})
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已注册!",
		})
		return
	}

	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致!",
		})
		c.Abort() // 终止处理链
		return
	}
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已注册!",
		})
		return
	}
	//判断phone是否重复
	phone := c.Query("phone")
	existingUserByPhone := models.FindUserByPhone(user.Phone)
	if existingUserByPhone.Phone != "" {
		c.JSON(-1, gin.H{
			"message": "该手机号已被注册!",
		})
		return
	}
	//判断email是否重复
	/*email := c.Query("email")
	existingUserByEmail := models.FindUserByEmail(user.Email)
	if existingUserByEmail.Email != "" {
		c.JSON(-1, gin.H{
			"message": "该邮箱已被注册!",
		})
		return
	}
	identity := c.Query("identity")
	existingUserByIdentity := models.FindUserByEmail(user.Identity)
	if existingUserByIdentity.Identity != "" {
		c.JSON(-1, gin.H{
			"message": "该身份重复注册!",
		})
		return
	}*/

	/*user.PassWord = password*/
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	/*user.Email = email*/
	user.Phone = phone
	/*user.Identity = identity*/
	user.LoginTime = time.Now()
	result := models.CreateUser(user) // 调用一次即可
	if result.Error != nil {
		c.JSON(-1, gin.H{
			"message": "新增用户失败：" + result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Email = c.PostForm("email")
	user.Phone = c.PostForm("phone")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //0正确 -1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, //0正确 -1失败
			"message": "修改用户成功",
			"data":    user,
		})
	}
}

// FindUserByNameAndPwd
// Summary 登陆
// @Tags 用户模块
// @param name query string false "name"
// @param password query string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	/*name := c.Query("name")
	password := c.Query("password")*/
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	clientIP := c.ClientIP()
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0正确 -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    0, //0正确 -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	user.ClientIp = clientIP
	models.UpdateUser(user)

	c.JSON(200, gin.H{
		"code":    0, //0正确 -1失败
		"message": "登陆成功",
		"data":    data,
	})
}

// LogoutUser
// @Summary 退出登录
// @Tags 用户模块
// @Param username query string true "用户名"
// @Success 200 {object} gin.H{"message": "退出登录成功"}
// @Failure 404 {object} gin.H{"message": "用户不存在"}
// @Router /user/logoutUser [post]
func LogoutUser(c *gin.Context) {
	username := c.Query("username") // 假设用户名存储在名为 "username" 的查询参数中
	user := models.FindUserByName(username)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"message": "用户不存在",
		})
		return
	}
	// 记录下线时间
	user.LoginOutTime = time.Now()
	// 计算登录时长
	loginDuration := time.Now().Sub(user.LoginTime)
	// 将登录时长赋值给 user.HeartbeatTime 字段
	user.HeartbeatTime = loginDuration
	user.IsLogout = true
	// 更新用户信息
	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"message": "退出登录成功",
	})
}

// 防止跨域的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	//将 HTTP 连接升级为 WebSocket 连接
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}
	//执行结束后关闭websocket连接，将资源释放
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	//基于tcp 转化为字节流  1表示文本类型
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(int(uint(id)))
	/*c.JSON(200, gin.H{
		"code":    0, //0正确 -1失败
		"message": "查询好友列表成功",
		"data":    users,
	})*/
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	/*targetID, _ := strconv.Atoi(c.Request.FormValue("targetId"))*/

	code, msg := models.AddFriend(uint(userId), targetName)
	/*c.JSON(200, gin.H{
		"code":    0, //0正确 -1失败
		"message": "查询好友列表成功",
		"data":    users,
	})*/
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	community := models.Community{}
	community.OwnerID = uint(ownerId)
	community.Name = name
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载群列表

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	//	name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
