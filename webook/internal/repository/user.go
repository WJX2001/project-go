package repository

import (
	"context"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 操作缓存
}

func (r *UserRepository) FindById(int64) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
}
