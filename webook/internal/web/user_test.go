package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"project-go/webook/internal/domain"
	"project-go/webook/internal/service"
	svcmocks "project-go/webook/internal/service/mocks"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
	}{}

	_, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(`
{
	"email": "123@qq.com",
	"password": "123456"
}
`)))
	require.NoError(t, err)

	// 这里可以继续使用 req
	//resp := httptest.NewRecorder()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 这里怎么拿到这个响应？
			handler := NewUserHandler(nil, nil)
			ctx := &gin.Context{}
			handler.SignUp(ctx)

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
	userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

	err := userSvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)

}

func TestUserHandler_SignUp1(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServiceInterface
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "1758807220@qq.com",
					Password: "Abcde1234@",
				}).Return(nil)
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234@",
    "passwordConfirm": "Abcde1234@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234@",
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)

				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@q",
    "password": "Abcde1234@",
    "passwordConfirm": "Abcde1234@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234@",
    "passwordConfirm": "Abcde1234@4"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234",
    "passwordConfirm": "Abcde1234"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "1758807220@qq.com",
					Password: "Abcde1234@",
				}).Return(service.ErrUserDuplicateEmail)
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234@",
    "passwordConfirm": "Abcde1234@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "1758807220@qq.com",
					Password: "Abcde1234@",
				}).Return(errors.New("随便一个error"))
				// 注册成功是return nil
				return userSvc
			},
			reqBody: `
	{
    "emailInfo": "1758807220@qq.com",
    "password": "Abcde1234@",
    "passwordConfirm": "Abcde1234@"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutesUser(server)

			req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.Use(func(c *gin.Context) {
				c.Set("user", UserClaims{})
			})
			t.Log(resp)
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}

}
