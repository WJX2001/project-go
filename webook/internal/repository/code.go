package repository

import (
	"context"
	"project-go/webook/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.CodeCache
}

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{cache: c}
}

func (repo *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
