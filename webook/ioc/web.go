package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"project-go/webook/internal/web"
	"project-go/webook/internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutesUser(server)
	//oauth2WechatHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		//logger.NewBuilder().AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder().
			JWTIgnorePaths("/user/signup").
			JWTIgnorePaths("/user/login").
			JWTIgnorePaths("/user/login_sms").
			JWTIgnorePaths("/oauth2/wechat/authurl").
			JWTIgnorePaths("/oauth2/wechat/callback").
			JWTIgnorePaths("/user/login_sms/code/send").
			Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}
