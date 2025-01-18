//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"project-go/webook/internal/repository"
	"project-go/webook/internal/repository/cache"
	"project-go/webook/internal/repository/dao"
	"project-go/webook/internal/service"
	"project-go/webook/internal/web"
	"project-go/webook/ioc"
)
import "github.com/google/wire"

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		InitDB, ioc.InitRedis,
		// 初始化 DAO
		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		// 直接基于内存实现
		ioc.InitSMSService,
		// 中间件呢？
		// 注册路由呢？
		web.NewUserHandler,
		//gin.Default,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
	//return gin.Default()
}
