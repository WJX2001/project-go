package service

import (
	"context"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里的问题
	// 将数据存起来
	return svc.repo.Create(ctx, u)
}
