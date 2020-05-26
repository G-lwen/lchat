package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
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
