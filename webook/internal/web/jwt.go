package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtHandler struct {
	// access_key
	atKey []byte
	// refresh_token key
	rtKey []byte
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"),
		rtKey: []byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"),
	}
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		// 设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"))
	if err != nil {
		return err
	}
	// 将token 放到header中
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid       int64
	UserAgent string
}

type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}
