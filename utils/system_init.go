package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig(hostIP string) {
	// 设置配置的名字是 app
	viper.SetConfigName("app")
	// 设置配置的路径
	viper.AddConfigPath("config")

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("读取配置文件错误:", err)
		return
	}
	mysqlDns := viper.GetString("mysql.dns")
	mysqlDns = fmt.Sprintf(mysqlDns, hostIP)
	fmt.Println("MySQ拼接字符串L", mysqlDns)
	// 打印配置内容
	fmt.Println("config app 里的内容", viper.Get("app"))

	// 初始化 MySQL 连接
	InitMySQL(mysqlDns)
}

func InitMySQL(dns string) {
	// 自定义日志，打印 SQL 语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{Logger: newLogger})
	if err != nil {
		// 处理数据库连接错误
		fmt.Println("初始化 MySQL 连接错误:", err)
		return
	}

	fmt.Println("数据库连接成功")
}

const (
	PublishKey = "websocket"
)

// ctx 上下文对象，可以用于在函数间传递取消信号、截止时间等
// channel 这是 Redis 的订阅频道名称，表示要发布消息到哪个频道
// msg 要发布的内容

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	//发布到指定频道
	fmt.Println("Publish...", msg)
	err = Red.Publish(ctx, channel, msg).Err() //发送消息
	return err
}

// Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	//创建订阅对象 订阅对象是一个在 Redis 客户端库中用于管理和处理订阅操作的实例
	sub := Red.Subscribe(ctx, channel)
	//接收订阅消息
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("Subscribe...", msg.Payload)
	return msg.Payload, err
}

// 与redis建立连接
/*func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	ctx := context.Background()
	pong, err := Red.Ping(ctx).Result()
	if err != nil {
		fmt.Println("初始化 Redis 失败", err)
	} else {
		fmt.Println("redis 初始化成功", pong)
	}

}*/
