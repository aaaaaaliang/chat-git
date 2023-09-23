package main

import (
	"chat/router"
	"chat/utils"
	"fmt"
	"os"
)

func main() {
	//初始化路径  数据库
	hostIP := os.Getenv("HOST_IP")
	// 初始化配置
	utils.InitConfig(hostIP)
	/*utils.InitRedis()*/
	r := router.Router()
	err := r.Run()
	if err != nil {
		fmt.Println("main.go", err)
		return
	}
}
