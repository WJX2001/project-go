package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicate
	ErrUserNotFound       = repository.ErrUserNotFound
)
var ErrInvalidUserOrPassword = errors.New("invalid user/password")

type UserServiceInterface interface {
	Login(ctx context.Context, u domain.User) (user domain.User, err error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context,
		user domain.User) error
	FindById(ctx context.Context,
		uid int64) (domain.User, error)
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserServiceInterface {
	return &UserService{
		repo: repo,
	}
}

// Login 登陆
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

// SignUp 注册
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

// UpdateNonSensitiveInfo 编辑信息
func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

// FindById 查找信息
func (svc *UserService) FindById(ctx context.Context,
	uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 这时候，怎么办
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断 有没有这个用户
	if err != repository.ErrUserNotFound {
		// 绝大部分请求会进来这里
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也会进来这里
		return u, err
	}
	// 明确知道，没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里会遇到主从延迟的问题
	//return u, err
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
}
