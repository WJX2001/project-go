package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
	ErrUserNotFound       = repository.ErrUserNotFound
)
var ErrInvalidUserOrPassword = errors.New("invalid user/password")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// 登陆
func (svc *UserService) Login(ctx context.Context, u domain.User) (user domain.User, err error) {
	// 先找用户
	findUser, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}
	return findUser, nil
}

// 注册
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

// 编辑信息
func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

// 查找信息
func (svc *UserService) FindById(ctx context.Context,
	uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}
