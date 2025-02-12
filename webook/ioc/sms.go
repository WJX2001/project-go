package ioc

import (
	"github.com/redis/go-redis/v9"
	"project-go/webook/internal/service/sms"
	"project-go/webook/internal/service/sms/memory"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	//// 换内存还是换别的
	//svc := memory.NewService()
	//
	//return ratelimit.NewRateLimitSMSService
	return memory.NewService()
}
