package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"project-go/webook/internal/service"
	"project-go/webook/internal/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登陆URL失败",
		})
		return
	}

	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
	//	State: state,
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		// 过期时间，你预期中一个用户完成登陆的时间
	//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
	//	},
	//})

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	// 校验一下state
	info, err := h.svc.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	//h.setJWTToken(ctx,)
	ctx.JSON(http.StatusOK, Result{
		Msg: "系统错误",
	})
}

//// OAuth2Handler 针对所有平台实现的方法
//type OAuth2Handler struct {
//	// 这里可以放不同的service
//	// wechatService
//	// dingdingService
//	// feishuService
//	svcs map[string]wechatService
//}

//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//	g := server.Group("/oauth2")
//	g.GET("/:platform/authurl", h.AuthURL)
//	g.Any("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(ctx *gin.Context) {
//	// 这里需要获取参数
//	platform := ctx.Param("platform")
//	//switch platform {
//	//case "wechat":
//	//
//	//}
//
//	svc := h.svcs[platform]
//
//}
//
//func (h *OAuth2Handler) Callback(ctx *gin.Context) {
//
//}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
