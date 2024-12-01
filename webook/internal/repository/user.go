package repository

import (
	"context"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository/cache"
	"project-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmailInfo
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 操作缓存
}

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

func (r *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

// 更新操作
func (r *UserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	return r.dao.UpdateById(ctx, r.toEntity(u))
}

func (r *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// TODO: 从数据库中取数据
	//u, err := r.dao.FindById(ctx, uid)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return r.toDomain(u), nil

	// TODO: 从缓存中取数据
	u, err := r.cache.Get(ctx, uid)
	if err == nil {
		// 必然有数据
		return u, nil
	}
	// 没这个数据
	//if err == cache.ErrKeyNotExist {
	//	// 去数据库里面加载
	//}

	// 如果Redis真的崩溃了，需要保护数据库
	// TODO:选加载: 数据库需要做限流
	// 如果选择不加载 用户体验不好

	ue, err := r.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志 做监控
		}
	}()

	return u, err
}
