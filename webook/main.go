package main

import (
	"github.com/gin-contrib/sessions/memstore"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"project-go/webook/internal/repository"
	"project-go/webook/internal/repository/dao"
	"project-go/webook/internal/service"
	user "project-go/webook/internal/web"
	"project-go/webook/internal/web/middleware"
	"project-go/webook/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	// 进行初始化
	//db := initDB()
	//server := initWebServer()
	//u := initUser(db)
	//u.RegisterRoutesUser(server)
	//server.Run("localhost:8080")

	// K8S部署web服务器，首先去除其他依赖(Mysql和Redis的干扰)
	server := gin.Default()
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world 你来了")
	})

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 使用中间件
	// 使用use 表明应用在server上的所有路由

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 使用第三方插件 通过redis 实现限流
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

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
	//	IgnorePaths("/user/login").
	//	Build())

	// JWT的中间件
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		JWTIgnorePaths("/user/signup").
		JWTIgnorePaths("/user/login").
		Build())

	return server
}

func initUser(db *gorm.DB) *user.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := user.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	// 进行初始化
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
