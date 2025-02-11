package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/repository"
	repomocks "project-go/webook/internal/repository/mocks"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		ctx      context.Context
		email    string
		password string

		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功", // 用户名和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Phone:    "13623231234",
					Password: "$2a$10$tc/oBMcdQIIe3OPotKYeE.sfkLCzbvYTwz4Fg7h2Mh8jbTGVU4wE.",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Phone:    "13623231234",
				Password: "$2a$10$tc/oBMcdQIIe3OPotKYeE.sfkLCzbvYTwz4Fg7h2Mh8jbTGVU4wE.",
				Ctime:    now,
			},

			wantErr: nil,
		},
		{
			name: "用户不存在", // 用户名和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对", // 用户名和密码是对的
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Phone:    "13623231234",
					Password: "$2a$10$tc/oBMcdQIIe3OPotKYeE.sfkLCzbvYTwz4Fg7h2Mh8jbTGVU4wE.",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体的测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(tc.ctx, domain.User{Email: tc.email, Password: tc.password})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
