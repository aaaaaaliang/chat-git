package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	// 创建Gin应用
	app := gin.Default()

	// 注册WebSocket路由
	app.GET("/ws", WebSocketHandler)

	// 启动应用
	err := app.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func WebSocketHandler(c *gin.Context) {
	// 获取WebSocket连接
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		panic(err)
	}

	// 处理WebSocket消息
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("messageType:", messageType)
		fmt.Println("p:", string(p))

		// 输出WebSocket消息内容
		c.Writer.Write(p)
	}

	// 关闭WebSocket连接
	ws.Close()
}
