package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	staffio "github.com/liut/staffio/client"
)

const (
	sKeyUser = "user"
)

type User = staffio.User

var (
	LoginHandler = gin.WrapF(staffio.LoginHandler)
	SetLoginPath = staffio.SetLoginPath
	SetAdminPath = staffio.SetAdminPath
)

func AuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := ginSession(c)
		if user, ok := sess.Get(sKeyUser).(*User); ok {
			if !user.IsExpired() {
				if user.NeedRefresh() {
					user.Refresh()
					sess.Set(sKeyUser, user)
					staffio.SessionSave(sess, c.Writer)
				}
				c.Set(sKeyUser, user)
				c.Next()
				return
			}
		}

		if redirect {
			c.Redirect(http.StatusFound, staffio.LoginPath)
			c.Abort()
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func UserWithContext(c *gin.Context) (user *User, ok bool) {
	v, ok := c.Get(sKeyUser)
	if ok {
		user = v.(*User)
	}
	if user == nil {
		log.Print("user not found in request")
	}

	return
}

// AuthCodeCallback Handler for Check auth with role[s] when auth-code callback
func AuthCodeCallback(roleName ...string) gin.HandlerFunc {
	return gin.WrapH(staffio.AuthCodeCallback(roleName...))
}

func HandlerShowMe(c *gin.Context) {
	user, ok := UserWithContext(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"me": user,
	})
}