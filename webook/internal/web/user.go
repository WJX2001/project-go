package user

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 在此定义跟 user有关的路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 不需要每次都编译，只需要暴露方法 进行预编译
func NewUserHandler() *UserHandler {
	// 定义正则表达式
	const (
		emailRegexPattern    = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
		passwordRegexPattern = `^(?=.*[0-9])(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

// 这种写法缺陷：容易被别人注册相同的路由
func (u *UserHandler) RegisterRoutesUser(server *gin.Engine) {
	// 统一处理前缀
	ug := server.Group("/user")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

// 注册
func (u *UserHandler) SignUp(ctx *gin.Context) {
	// 定义在里面 防止其他人调用
	type SignUpReq struct {
		Email           string `json:"emailInfo"`
		PasswordConfirm string `json: "passwordConfirm"`
		Password        string `json: "password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写会一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 判断邮箱格式
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		fmt.Println(err)
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if req.PasswordConfirm != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	ctx.String(http.StatusOK, "hello 你在注册")
	fmt.Printf("%v", req)
	// 这边就是数据库操作
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
