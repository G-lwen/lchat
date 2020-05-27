package router

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"lchat/service/entity"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": err.Error(),
		})
		return
	}

    now := time.Now()
	filepath := fmt.Sprintf("data/images/%04d/%02d/%02d/%s.%s",
		now.Year, now.Month(), now.Day(), uuid.NewV4().String(), file.Filename[strings.LastIndex(file.Filename, "."):])

	if err = c.SaveUploadedFile(file, fmt.Sprintf("data/images/%s", filepath)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("/images/%s", filepath),
	})
}

func downloadPost(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	u, _ := user.(*entity.User);
	postId := c.Param("postId")
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
	if err = post.Load(); err != nil || (!post.Published && post.UserId != u.ID) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"message": err.Error(),
		})
		return
	}
	c.Header("Content-Type", "application/force-download")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s%s", post.Title, ".md"))
	c.Header("Content-Transfer-Encoding", "binary")
	w := c.Writer
	w.WriteHeader(http.StatusOK)
	w.WriteString(post.Body)
}