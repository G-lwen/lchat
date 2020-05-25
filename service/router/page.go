package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"lchat/service/entity"
	"net/http"
	"strconv"
)

// 首页
func indexPage(c *gin.Context) {
	var (
		isAdmin bool
		isLogin bool
		nickName string
		countTags int
	)
	pageIndex := c.DefaultQuery("pageIndex", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		isLogin = true
		nickName = u.NickName
		if u.IsAdmin {
			isAdmin = true
		}
	}
	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 1 {
		index = 1
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil || size < 1 {
		size = 10
	}
	posts, count, _ := entity.ListPosts(index, size)
	for _, post := range posts {
		post.LoadTags()
	}
	tags, _ := entity.ListAllTags()
	for _, tag := range tags {
		countTags += tag.Total
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"isLogin": isLogin,
		"isAdmin": isAdmin,
		"nickName": nickName,
		"pageIndex": index,
		"pageSize": size,
		"posts": posts,
		"totalPage": count / size + 1,
		"tags": tags,
		"countTags": countTags,
	})
}

// 登录页面
func loginPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if _, ok := user.(*entity.User); ok {
		c.Redirect(http.StatusSeeOther, "/")
	}
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

// 注册页面
func registerPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if _, ok := user.(*entity.User); ok {
		c.Redirect(http.StatusSeeOther, "/")
	}
	c.HTML(http.StatusOK, "register.tmpl", gin.H{})
}

// 写博客页面
func postEditPage(c *gin.Context) {
	var (
		isAdmin bool
		nickName string
		hasPost bool
		post = &entity.Post{}
	)
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		nickName = u.NickName
		if u.IsAdmin {
			isAdmin = true
		}
	} else {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
    postId := c.DefaultQuery("postId", "0")
    pid, err := strconv.ParseUint(postId, 10, 64)
    if err == nil && pid > 0 {
    	post.ID = uint(pid)
    	if err = post.Load(); err == nil {
    		hasPost = true
    		post.LoadTags()
		}
	}
	c.HTML(http.StatusOK, "post_edit.tmpl", gin.H{
		"isAdmin": isAdmin,
		"nickName": nickName,
		"hasPost": hasPost,
		"post": post,
	})
}

// ping 服务
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 文章页面
func postPage(c *gin.Context) {
	var (
		isAdmin bool
		isLogin bool
		nickName string
		countTags int
		userId uint
		post = &entity.Post{}
	)
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		isLogin = true
		userId = u.ID
		nickName = u.NickName
		if u.IsAdmin {
			isAdmin = true
		}
	}
	postId := c.Query("postId")
	pid, _ := strconv.ParseUint(postId, 10, 64)
	post.ID = uint(pid)
	if err := post.Load(); err != nil || (!post.Published && post.UserId != userId) {
		post.Summary = ""
		post.Body = ""
		post.HtmlBody = "<h1>没有该文章的访问权限</h1>"
	} else {
		post.LoadTags()
	}
	tags, _ := entity.ListAllTags()
	for _, tag := range tags {
		countTags += tag.Total
	}
	c.HTML(http.StatusOK, "post.tmpl", gin.H{
		"isLogin": isLogin,
		"isAdmin": isAdmin,
		"nickName": nickName,
		"tags": tags,
		"countTags": countTags,
		"post": post,
	})
}

//用户文库
func userPostsPage(c *gin.Context) {
	var (
		isAdmin bool
		nickName string
		userId uint
	)
	pageIndex := c.DefaultQuery("pageIndex", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		userId = u.ID
		nickName = u.NickName
		if u.IsAdmin {
			isAdmin = true
		}
	} else {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	index, err := strconv.Atoi(pageIndex)
	if err != nil || index < 1 {
		index = 1
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil || size < 1 {
		size = 10
	}
	posts, count, _ := entity.ListPostsByUserId(index, size, userId, true)
	for _, post := range posts {
		if len(post.Summary) > 50 {
			post.Summary = post.Summary[:50] + "..."
		}
		post.LoadTags()
	}
	c.HTML(http.StatusOK, "post_list.tmpl", gin.H{
		"isAdmin": isAdmin,
		"nickName": nickName,
		"pageIndex": index,
		"pageSize": size,
		"posts": posts,
		"totalPage": count / size + 1,
	})
}
