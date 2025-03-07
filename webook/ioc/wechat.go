package ioc

import (
	"github.com/google/uuid"
	//"os"
	"project-go/webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("么有找到环境变量 WECHAT_APP_ID")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("没有找到环境变量 WECHAT_APP_SECRET")
	//}
	// 自己造一个
	appId := uuid.New().String()
	appKey := uuid.New().String()

	return wechat.NewService(appId, appKey)
}
