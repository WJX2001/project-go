package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 如何判断，session过期时间到了
		updateTime := sess.Get("update_Time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		// 说明还没有刷新过(刚登陆，还没刷新过)
		if updateTime == nil {
			sess.Set("update_Time", now)
			sess.Save()
			return
		}

		// 如果有updateTime
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_Time", now)
			sess.Save()
			return
		}

	}
}
