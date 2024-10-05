package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

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

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	// 将数据存起来
	return svc.repo.Create(ctx, u)
}
