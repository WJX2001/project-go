package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	// 编译器会在编译的时候，把set_code的代码 放进来这个 luaSetCode 变量里
	//go:embed lua/set_code.lua
	luaSetCode         string
	ErrCodeSendTooMany = errors.New("发送太频繁")
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	case -2:
		// 系统错误
		return errors.New("验证码存在，但是没有过期时间")
	default:
		return nil
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
