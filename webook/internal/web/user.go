package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-go/webook/internal/domain"
	service "project-go/webook/internal/service"
	"time"
)

const biz = "login"

// 确保 UserHandler上实现了 handler接口
//var _ handler = (*UserHandler)(nil)

// UserHandler 在此定义跟 user有关的路由
type UserHandler struct {
	svc         service.UserServiceInterface
	codeSvc     service.CodeServiceInterface
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	jwtHandler
}

// 不需要每次都编译，只需要暴露方法 进行预编译
func NewUserHandler(svc service.UserServiceInterface, codeSvc service.CodeServiceInterface) *UserHandler {
	// 定义正则表达式
	const (
		emailRegexPattern    = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
		passwordRegexPattern = `^(?=.*[0-9])(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,72}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
	}
}

// 这种写法缺陷：容易被别人注册相同的路由
func (u *UserHandler) RegisterRoutesUser(server *gin.Engine) {
	// 统一处理前缀
	ug := server.Group("/user")
	ug.GET("/profile", u.ProfileByJWT)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/logout", u.Logout)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 这里也可以加上各种校验
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办？
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码校验通过",
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 是不是一个合法的手机号码
	// 考虑用正则表达式去校验
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	//if err != nil {
	//	//ctx.JSON(http.StatusOK, "系统异常")
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg:  "系统错误",
	//	})
	//	return
	//}
	//ctx.JSON(http.StatusOK, Result{
	//	Msg: "发送成功",
	//})
}

// 注册
func (u *UserHandler) SignUp(ctx *gin.Context) {
	// 定义在里面 防止其他人调用
	type SignUpReq struct {
		Email           string `json:"emailInfo"`
		ConfirmPassword string `json:"passwordConfirm"`
		Password        string `json:"password"`
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

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	// 调用一下 svc 的方法 进行注册
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	return

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		PassWord string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 svc
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.PassWord,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// TODO: 在这里使用JWT登陆态
	// 生成一个 JWT token
	// 后续要在JWT token里面带上个人数据
	// 比如：userID
	//token := jwt.New(jwt.SigningMethodHS512)

	//if done {
	//	return
	//}
	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登陆成功")
	return
}

//func (u *UserHandler) setJWTTOKEN(ctx *gin.Context, uid int64) error {
//	claims := UserClaims{
//		// 设置过期时间
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
//		},
//		Uid:       uid,
//		UserAgent: ctx.Request.UserAgent(),
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenStr, err := token.SignedString([]byte("IjkxUQzY7dMQ4gdYLUMVvMXsIpl1E7f4"))
//	if err != nil {
//		return err
//	}
//	// 将token 放到header中
//	ctx.Header("x-jwt-token", tokenStr)
//	return nil
//}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		PassWord string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 svc
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.PassWord,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 在这里登陆成功了
	// 将session取出来
	sessionInfo := sessions.Default(ctx)
	// 可以随便设置，放入session的值
	sessionInfo.Set("userId", user.Id)
	sessionInfo.Save() // sessionInfo进行保存
	ctx.String(http.StatusOK, "登陆成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登陆成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 改邮箱 密码 或者能不能改手机号
		Nickname string `json:"nickname"`
		// YYYY-MM-DD
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 从session中获取登陆状态
	sess := sessions.Default(ctx)

	sessionId := sess.Get("userId")

	// 进行数据操作
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}

	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       sessionId.(int64),
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "更新成功")

}

func (u *UserHandler) ProfileByJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	// 可以断定 必然有 claims
	//if !ok {
	//	// 可以考虑监控这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	//
	// ok 代表是不是 *UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 补充其他代码
	userInfo, err := u.svc.FindById(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname: userInfo.Nickname,
		Email:    userInfo.Email,
		AboutMe:  userInfo.AboutMe,
		Birthday: userInfo.Birthday.Format(time.DateOnly),
	})

}

func (u *UserHandler) Profile(ctx *gin.Context) {

	// 从session中获取登陆状态
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	//userInfo, err := u.svc.FindById(ctx, sessionId)
	userInfo, err := u.svc.Profile(ctx, id)

	// TODO: 通过JWT token中的UID获取用户信息

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname: userInfo.Nickname,
		Email:    userInfo.Email,
		AboutMe:  userInfo.AboutMe,
		Birthday: userInfo.Birthday.Format(time.DateOnly),
	})

}

//type UserClaims struct {
//	jwt.RegisteredClaims
//	// 声明你自己的要放进去 token 里面的数据
//	Uid       int64
//	UserAgent string
//	// 自己随便添加 想要添加的属性
//	// 密码和权限等敏感数据 不要放在这里
//}
