package ratelimit

import (
	"context"
	"fmt"
	"project-go/webook/internal/service/sms"
	"project-go/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRateLimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		/*
			1. 可以限流：保守策略，你的下游很坑的时候
			2. 可以不限：你的下游很强，业务可用性要求很高，尽量容错策略
			3. 包一下这个错误
		*/
		return err
	}

	if limited {
		return errLimited
	}

	// 在这里加一些代码，新特性
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 在这里也可以加一些代码，新特性
	return err
}
