package repository

import (
	"context"
	"project-go/webook/internal/repository/cache"
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}
type CacheCodeRepository struct {
	cache cache.CodeCache
}

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{cache: c}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CacheCodeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
