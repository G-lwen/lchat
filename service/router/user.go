package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"lchat/service/entity"
	"lchat/service/utils"
	"net/http"
	"strings"
)

// 获取邮箱注册码
func getRegisterCode(c *gin.Context) {
	email := c.Query("email")
	if !utils.VerifyEmailFormat(email) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "邮箱格式错误",
		})
		return
	}

	code := utils.RandomCodeGenerator(6, utils.NumberCharasterCode)
	err := utils.SendEmailRegisterCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	} else {
		session := sessions.Default(c)
		session.Set(email, code)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "发送成功，请注意接收!",
		})
	}
}

// 用户注册
func register(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	code := c.PostForm("code")

	session := sessions.Default(c)
	realCode := session.Get(email)
	if rc, ok := realCode.(string); ok && code != "" && strings.ToLower(code) != strings.ToLower(rc) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "邮箱验证码错误!",
		})
		return
	}

	user := &entity.User{
		Email: email,
		Password: password,
		NickName: utils.RandomCodeGenerator(12, utils.NumberLowerCharasterCode),
	}

	user.LoadByEmail()
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "邮箱号已被注册",
		})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "注册成功，请前往登录!",
		})
	}
}

// 用户登录
func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &entity.User{
		Email: username,
	}

	if err := user.LoadByEmail(); err == nil && user.VerificationPassword(password) {
		session := sessions.Default(c)
		session.Clear()
		session.Set(sessionKey, user)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "ok",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 401,
		"message": "用户帐号或密码错误",
	})
}

// 用户登出
func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
