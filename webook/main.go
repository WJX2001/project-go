package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project-go/webook/internal/repository"
	"project-go/webook/internal/repository/dao"
	"project-go/webook/internal/service"
	user "project-go/webook/internal/web"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 进行初始化
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutesUser(server)
	server.Run("localhost:8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// 使用中间件
	// 使用use 表明应用在server上的所有路由
	server.Use(cors.New(cors.Config{
		// AllowOrigins:  []string{"http://localhost:8000"},
		AllowMethods:  []string{"POST", "GET"},
		AllowHeaders:  []string{"content-type", "Authorization"},
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
