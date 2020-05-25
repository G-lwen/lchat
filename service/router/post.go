package router

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"lchat/service/entity"
	"net/http"
	"strconv"
	"time"
)

// 文章发布
func postPublish(c *gin.Context) {
	var (
		userId uint
		post *entity.Post
	)
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		userId = u.ID
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"message": "没有权限访问",
		})
		return
	}
	postId := c.DefaultPostForm("postId", "0")
	title := c.PostForm("title")
	body := c.PostForm("body")
	htmlBody := c.PostForm("htmlBody")

	pid, err := strconv.ParseUint(postId, 10, 64)
	if err == nil && pid > 0 {
		post = &entity.Post{}
		post.ID = uint(pid)
		if err = post.Load(); err == nil && post.UserId == userId {
			if title != "" {
				post.Title = title
			}
			post.Body = body
			post.HtmlBody = htmlBody
			post.ExtractSummary()
			if err = post.Update(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 500,
					"message": "系统繁忙，请稍后重试!",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"message": "ok",
					"postId": post.ID,
				})
			}
			return
		}
	}
	if title == "" {
		now := time.Now()
		title = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}
	post = &entity.Post{
		UserId: userId,
		Title: title,
		Body: body,
		HtmlBody: htmlBody,
		Published: true,
	}
	post.ExtractSummary()
	if err := post.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "ok",
		"postId": post.ID,
	})
}


// 文章添加标签
func postAddTag(c *gin.Context) {
	var userId uint
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		userId = u.ID
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"message": "没有权限访问",
		})
		return
	}
	postId := c.PostForm("postId")
	tagName := c.PostForm("tagName")

	if tagName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "空标签错误, 添加失败",
		})
		return
	}
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": err.Error(),
		})
		return
	}
	post := &entity.Post{}
	post.ID = uint(pid)
	if err = post.Load(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	if userId != post.UserId {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "非法操作",
		})
		return
	}
	tag := &entity.Tag{
		Name: tagName,
	}
	if err = tag.LoadByName(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	if err = post.AddTag(tag.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "ok",
	})
}

// 文章移除标签
func postRemoveTag(c *gin.Context) {
	var userId uint
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	if u, ok := user.(*entity.User); ok {
		userId = u.ID
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"message": "没有权限访问",
		})
		return
	}
	postId := c.Query("postId")
	tagName := c.Query("tagName")

	if tagName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "空标签错误, 删除失败",
		})
		return
	}
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": err.Error(),
		})
		return
	}
	post := &entity.Post{}
	post.ID = uint(pid)
	if err = post.Load(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	if userId != post.UserId {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": "非法操作",
		})
		return
	}
	tag := &entity.Tag{
		Name: tagName,
	}
	if err = tag.LoadByName(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	if err = post.RemoveTag(tag.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "ok",
	})
}

