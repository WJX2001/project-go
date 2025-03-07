//go:build manual

package wechat

import (
	"os"
	"testing"
)

// 手动验证 提前验证代码
func Test_service_VerifyCode(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID")
	}

	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	svc := NewService(appId, appKey)
	//svc.VerifyCode(context.Background(), "")
}
