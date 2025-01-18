package repository

import (
	"context"
	"database/sql"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository/cache"
	"project-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, u domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	//return domain.User{
	//	Id:       u.Id,
	//	Email:    u.Email,
	//	Password: u.Password,
	//}, nil
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	//return r.dao.Insert(ctx, dao.User{
	//	Email:    u.Email,
	//	Password: u.Password,
	//})
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

//func (r *CachedUserRepository) toDomain(u dao.User) domain.User {
//	return domain.User{
//		Id:       u.Id,
//		Email:    u.Email,
//		Password: u.Password,
//		AboutMe:  u.AboutMe,
//		Nickname: u.Nickname,
//		Birthday: time.UnixMilli(u.Birthday),
//	}
//}

func (r *CachedUserRepository) toEntity(u domain.User) dao.User {
	//return dao.User{
	//	Id:       u.Id,
	//	Email:    u.Email,
	//	Password: u.Password,
	//	Birthday: u.Birthday.UnixMilli(),
	//	AboutMe:  u.AboutMe,
	//	Nickname: u.Nickname,
	//}
	return r.domainToEntity(u)
}

// 更新操作
func (r *CachedUserRepository) UpdateNonZeroFields(ctx context.Context, u domain.User) error {
	return r.dao.UpdateById(ctx, r.toEntity(u))
}

func (r *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
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

	//u = domain.User{
	//	Id:       ue.Id,
	//	Email:    ue.Email,
	//	Password: ue.Password,
	//}
	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志 做监控
		}
	}()

	return u, err
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
