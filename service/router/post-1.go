package router
//
//import (
//	jwt "github.com/appleboy/gin-jwt/v2"
//	"github.com/gin-gonic/gin"
//	"lchat/service/entity"
//	"net/http"
//	"strconv"
//)
//
//// 创建文章
//func PostCreate(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	title := c.PostForm("title")
//	body := c.PostForm("body")
//
//	if title == "" {
//		title = "Untitled"
//	}
//	post := &entity.Post{
//		UserId: user.(*entity.User).ID,
//		Title: title,
//		Body: body,
//	}
//	if err := post.Save(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//		"post": post,
//	})
//}
//
//// 编辑文章
//func PostEdit(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//	title := c.PostForm("title")
//	body := c.PostForm("body")
//
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := new(entity.Post)
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "无法找到对应的文章",
//		})
//		return
//	}
//	if post.UserId != user.(*entity.User).ID {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非文章作者无法进行编辑",
//		})
//		return
//	}
//	post.Title = title
//	post.Body = body
//	if err = post.Update(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	post.LoadTags()
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//		"post": post,
//	})
//}
//
//// 发布或者取消发布文章
//func PostPublish(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := &entity.Post{}
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if user.(*entity.User).ID != post.UserId {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非法操作",
//		})
//		return
//	}
//	post.Published = !post.Published
//	if err = post.Update(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//	})
//}
//
//// 根据 ID 查看文章
//func PostGet(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := &entity.Post{}
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if !post.Published && user.(*entity.User).ID != post.UserId {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非法浏览文章内容",
//		})
//		return
//	}
//	post.LoadTags()
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//		"post": post,
//	})
//}
//
//// 根据 ID 删除文章
//func PostDelete(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := &entity.Post{}
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "无法找到对应的文章",
//		})
//		return
//	}
//	if user.(*entity.User).ID != post.UserId {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非法删除文章",
//		})
//		return
//	}
//	if err = post.Delete(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//	})
//}
//
//// 获取文章列表
//func PostsGet(c *gin.Context) {
//	var (
//		index, size int
//		err error
//		posts []*entity.Post
//		count uint
//	)
//	pageIndex := c.DefaultQuery("pageIndex", "1")
//	pageSize := c.DefaultQuery("pageSize", "10")
//	tagName := c.Query("tagName")
//	index, err = strconv.Atoi(pageIndex)
//	if err != nil || index < 1 {
//		index = 1
//	}
//	size, err = strconv.Atoi(pageSize)
//	if err != nil || size < 1 {
//		size = 10
//	}
//	if tagName == "" {
//		if posts, err = entity.ListPosts(index, size); err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 500,
//				"message": err.Error(),
//			})
//			return
//		}
//		count = entity.CountPosts()
//	} else {
//		tag := &entity.Tag{
//			Name: tagName,
//		}
//		if err = tag.LoadByName(); err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 500,
//				"message": err.Error(),
//			})
//			return
//		}
//		if posts, err = entity.ListPostsByTag(index, size, tag.ID); err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 500,
//				"message": err.Error(),
//			})
//			return
//		}
//		count = entity.CountPostsByTag(tag.ID)
//	}
//	for _, post := range posts {
//		post.ExtractSummary()
//		post.LoadTags()
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//		"posts": posts,
//		"total": count,
//	})
//}
//
//// 获取个人的文章列表
//func PersonalPostsGet(c *gin.Context) {
//	var (
//		index, size int
//		err error
//		posts []*entity.Post
//		count uint
//		isPublish bool
//	)
//	user, _ := c.Get(jwt.IdentityKey)
//	pageIndex := c.DefaultQuery("pageIndex", "1")
//	pageSize := c.DefaultQuery("pageSize", "10")
//	published := c.DefaultQuery("published", "true")
//	index, err = strconv.Atoi(pageIndex)
//	if err != nil || index < 1 {
//		index = 1
//	}
//	size, err = strconv.Atoi(pageSize)
//	isPublish, err = strconv.ParseBool(published)
//	if err != nil {
//		isPublish = true
//	}
//	if err != nil || size < 1 {
//		size = 10
//	}
//	if posts, err = entity.ListPostsByUserId(index, size, user.(*entity.User).ID, isPublish); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	count = entity.CountPosts()
//	for _, post := range posts {
//		post.ExtractSummary()
//		post.LoadTags()
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//		"posts": posts,
//		"total": count,
//	})
//}
//
//// 获取标签列表
//func TagsGet(c *gin.Context) {
//	tags, err := entity.ListAllTags()
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code":    200,
//		"message": "ok",
//		"tags":    tags,
//		"total":   entity.CountTags(),
//	})
//}
//
//// 文章添加标签
//func PostAddTag(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//	tagName := c.PostForm("tagName")
//
//	if tagName == "" {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "空标签错误, 添加失败",
//		})
//		return
//	}
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := &entity.Post{}
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if user.(*entity.User).ID != post.UserId {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非法操作",
//		})
//		return
//	}
//	tag := &entity.Tag{
//		Name: tagName,
//	}
//	if err = tag.LoadByName(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if err = post.AddTag(tag.ID); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//	})
//}
//
//// 文章移除标签
//func PostRemoveTag(c *gin.Context) {
//	user, _ := c.Get(jwt.IdentityKey)
//	postId := c.Param("postId")
//	tagName := c.Query("tagName")
//
//	if tagName == "" {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "空标签错误, 删除失败",
//		})
//		return
//	}
//	pid, err := strconv.ParseUint(postId, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": err.Error(),
//		})
//		return
//	}
//	post := &entity.Post{}
//	post.ID = uint(pid)
//	if err = post.Load(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if user.(*entity.User).ID != post.UserId {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 404,
//			"message": "非法操作",
//		})
//		return
//	}
//	tag := &entity.Tag{
//		Name: tagName,
//	}
//	if err = tag.LoadByName(); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	if err = post.RemoveTag(tag.ID); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 500,
//			"message": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"message": "ok",
//	})
//}
