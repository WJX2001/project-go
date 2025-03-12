package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"project-go/webook/pkg/logger"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	logger        logger.LoggerV1
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: func(ctx context.Context, al *AccessLog) {

		},
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			// Body 读完就没有了
			body, _ := io.ReadAll(ctx.Request.Body)
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			ctx.Request.GetBody = func() (io.ReadCloser, error) {
				return reader, nil
			}

			if len(body) > 1024 {
				body = body[:1024]
			}
			// 这其实是一个很消耗 CPU 和 内存的操作
			// 因为会引起复制
			al.RespBody = string(body)
		}
		defer func() {
			//duration := time.Since(start)
			if b.allowReqBody && ctx.Request.Body != nil {

			}
			b.loggerFunc(ctx, al)
		}()

		// 执行到业务逻辑
		ctx.Next()

	}
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求 URL
	Url      string
	Duration time.Duration
	ReqBody  string
	RespBody string
}
