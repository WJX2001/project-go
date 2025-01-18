package ioc

import (
	"project-go/webook/internal/service/sms"
	"project-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 换内存还是换别的
	return memory.NewService()
}
