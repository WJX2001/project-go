package main

import (
	user "project-go/webook/internal/web"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	u := user.NewUserHandler()
	u.RegisterRoutesUser(server)
	server.Run("localhost:8080")
}
