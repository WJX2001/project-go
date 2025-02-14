package retryable

import (
	"context"
	"errors"
	"project-go/webook/internal/service/sms"
)

// 小心并发问题
type service struct {
	svc sms.Service
	// 重试
	retryCnt int
}

func (s service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	for err != nil && s.retryCnt < 10 {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}
