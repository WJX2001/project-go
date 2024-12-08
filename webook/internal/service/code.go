package service

import (
	"context"
	"fmt"
	"math/rand"
	"project-go/webook/internal/repository"
	"project-go/webook/internal/service/sms"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

const codeTplId = "1877556"

// Send 发验证码 我需要什么参数
func (svc *CodeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	code string,
	phone string,
) error {
	// 三个步骤：
	// 生成一个验证码
	code1 := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code1)
	if err != nil {
		// 有问题
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code1}, phone)
	return err
}

func (svc *CodeService) generateCode() string {
	// 六位数 num 在 0, 999999 之间， 闭区间
	num := rand.Intn(1000000)
	// 不够六位的，加上前导0
	return fmt.Sprintf("%6d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string,
) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

//
//func (svc *CodeService) VerifyV1(ctx context.Context, biz string,
//	phone string, inputCode string,
//) error {
//
//}
