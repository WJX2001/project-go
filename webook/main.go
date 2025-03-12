package main

import (
	"bytes"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"project-go/webook/internal/web/middleware"
	"strings"
	"time"
)

func main() {
	//initViper()
	//initViperV1()
	//initViperRemote()
	initLogger()
	// TODO: 使用 wire进行改造
	server := InitWebServer()
	server.Run(":8082")
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := ``
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 使用中间件
	// 使用use 表明应用在server上的所有路由

	//redisClient := redis.NewClient(&redis.Options{
	//	//Addr: "localhost:6379",
	//	// 这里需要连接到K8s部署的redis
	//	//Addr: "webook-live-redis:11479",
	//	// 直接使用配置文件中的
	//	Addr: config.Config.Redis.Addr,
	//})

	// 使用第三方插件 通过redis 实现限流
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins:  []string{"http://localhost:8000"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"content-type", "Authorization"},
		// 不加此配置，前端拿不到jwt-token
		ExposeHeaders: []string{"x-jwt-token"}, // 后续JWT会使用
		// 是否允许带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// return origin == "https://github.com"
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 配置session
	/**
	步骤一（选择存储方式）：
		同理也可换成 memStore 和 Redis
	*/
	//store := cookie.NewStore([]byte("secret"))
	//store.Options(sessions.Options{
	//	Secure:   true,  // 仅通过HTTPS传输Cookie
	//	HttpOnly: true,  // 禁止JavaScript访问Cookie
	//})

	// TODO: 替换存储为redis
	/**
		第一个参数是最大空闲连接数量
		第二个就是 tcp
	    第三个，第四个就是连接信息和密码
	    第五，第六就是两个key
	*/
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	// authentication key, encryption key
	//	/**
	//		authentication: 是指身份认证
	//	    encryption: 是指数据加密
	//	    这两者再加上授权（权限控制），就是信息安全的三个核心概念
	//	*/
	//	[]byte("fdDxNKZ6hNsXe1Ax5GWjbSlTKNhxSmZU"),
	//	[]byte("rcziTpeJ0dhwGKN6v3sHBCu92J0pmK9y"))
	//if err != nil {
	//	panic(err)
	//}

	// TODO: store替换成 memCache
	store := memstore.NewStore([]byte("fdDxNKZ6hNsXe1Ax5GWjbSlTKNhxSmZU"), []byte("rcziTpeJ0dhwGKN6v3sHBCu92J0pmK9y"))
	server.Use(sessions.Sessions("mysession", store))
	// 步骤三

	// Session 的中间件
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/user/signup").
	//	IgnorePaths("/user/login_sms/code/send").
	//	IgnorePaths("/user/login_sms").
	//	IgnorePaths("/user/login").
	//	Build())

	// JWT的中间件
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		JWTIgnorePaths("/user/signup").
		JWTIgnorePaths("/user/login").
		JWTIgnorePaths("/user/login_sms").
		JWTIgnorePaths("/user/login_sms/code/send").
		Build())

	return server
}

func initViperV1() {
	viper.SetDefault("db.mysql.dsn",
		"root:root@tcp(localhost:3306)/mysql")
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initViper() {
	viper.SetDefault("db.mysql.dsn",
		"root:root@tcp(localhost:3306)/mysql")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go .yaml 之类的后缀名
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是yaml格式
	// 现实有很多格式，JSON,XML,YAML
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	// 读取配置到 viper里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3",
		// 通过webook 和其他使用 etcd的区别出来
		"127.0.0.1:12379", "")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 如果不 replace,直接用 zap.L() 你啥都打不出来
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")

	type Demo struct {
		Name string `json:"name"`
	}

	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}
