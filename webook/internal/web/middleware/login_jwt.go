package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"project-go/webook/internal/web"
	"strings"
	"time"
)

// LoginJWTMiddlewareBuilder JWT 登陆校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func (l *LoginJWTMiddlewareBuilder) JWTIgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 现在使用JWT来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			// 没登陆 有人乱传 token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		// 使用指针的原因
		/**
		因为ParseWithClaims 将会修改claims中的数据，如果不传入指针，相当于复制了一份，这并没有什么用处
		ParseWithClaims 里面一定要传入指针
		*/
		claims := &web.UserClaims{}
		//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"), nil
		//})
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"), nil
		})

		if err != nil {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err 为 nil, token 不为nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
