package router

import (
	"encoding/gob"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"lchat/service/conf"
	"lchat/service/entity"
	"lchat/service/middleware"
	"lchat/service/utils"
	"net/http"
	"os"
	"time"
)

const (
	sessionKey = "lchat-session-key"
)

func Run() error {

	setWebLog()

	router := gin.Default()

	setSession(router)
	setView(router)

	// 允许跨域请求
	router.Use(middleware.Cors())

	router.GET("/ping", ping)
	router.GET("/", indexPage)
	router.GET("/login", loginPage)
	router.POST("/login", login)
	router.GET("/logout", logout)
	router.GET("/register", registerPage)
	router.GET("/registerCode", getRegisterCode)
	router.POST("/register", register)
	router.GET("/post", postPage)
	router.GET("/oauth/:oauthType", oauth)
	router.GET("/oauth/:oauthType/callback", oauthCallback)

	authorized := router.Group("/")
	authorized.Use(auth())
	{
		authorized.GET("/post/edit", postEditPage)
		authorized.POST("/post/publish", postPublish)
		authorized.POST("/post/addTag", postAddTag)
		authorized.DELETE("/post/removeTag", postRemoveTag)
		authorized.GET("/user/posts", userPostsPage)
		authorized.POST("/upload", upload)
	}

	return router.Run(":" + conf.Get().Server.Port)
}

//设置 Gin 日志
func setWebLog() {
	gin.DisableConsoleColor()
	f, _ := os.Create("data/logs/web.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
}

// 设置Session
func setSession(router *gin.Engine) {
	// session保存负责的类型需要先进行注册
	gob.Register(&entity.User{})

	store := sessions.NewCookieStore([]byte("lchat-session-secret"))
	router.Use(sessions.Sessions("lchat-session", store))
}

// 设置视图
func setView(router *gin.Engine) {
	router.Static("/css", "ui/static/css")
	router.Static("/js", "ui/static/js")
	router.Static("/img", "ui/static/img")
	router.Static("/fonts", "ui/static/fonts")
	router.Static("/images", "data/images")
	router.StaticFile("/favicon.ico", "ui/templates/favicon.ico")
	//router.LoadHTMLGlob("ui/templates/**/*")
	funcMap := template.FuncMap{
		"timeFormat": func(data time.Time, layout string) string {
			return data.Format(layout)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"minus": func(a, b int) int {
			return a - b
		},
	}
	router.SetFuncMap(funcMap)
	router.LoadHTMLFiles("ui/templates/index.tmpl",
		"ui/templates/user/login.tmpl", "ui/templates/user/register.tmpl",
		"ui/templates/post/post_edit.tmpl", "ui/templates/post/post.tmpl",
		"ui/templates/post/post_list.tmpl",
	)
}

// 权限认证
func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if user := session.Get(sessionKey); user != nil {
			if u, ok := user.(*entity.User); ok && (!utils.URLPathMatch("/admin/**", c.Request.URL.Path) || u.IsAdmin) {
				c.Next()
				return
			} else {
				c.HTML(http.StatusForbidden, "errors/errors.html", gin.H{
					"message": "该用户没有访问权限",
				})
				c.Abort()
			}
		} else {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
		}
	}
}