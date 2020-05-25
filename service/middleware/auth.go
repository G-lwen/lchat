package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"lchat/service/entity"
	"lchat/service/utils"
	"time"
)

// user jwt
func AuthRequired() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "lchat-realm",
		Key:         []byte("lchat-secret"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: jwt.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entity.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			user := &entity.User{
				Email: claims[jwt.IdentityKey].(string),
			}
			user.LoadByEmail()
			return user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			username := c.PostForm("username")
			password := c.PostForm("password")

			user := &entity.User{
				Email: username,
			}
			err := user.LoadByEmail()
			if err == nil && user.VerificationPassword(password) {
				return user, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if user, ok := data.(*entity.User); ok &&
				(!utils.URLPathMatch("/admin/**", c.Request.URL.Path) || user.IsAdmin) {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	})
}
