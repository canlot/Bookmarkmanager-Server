package Handlers

import (
	"Bookmarkmanager-Server/Configuration"
	"Bookmarkmanager-Server/Models"
	"github.com/akyoto/cache"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

var tokenCache *cache.Cache

func init() {
	tokenCache = cache.New(1 * time.Hour)
}

func Authenticate(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	if strings.Split(bearerToken, " ")[0] != "Bearer" {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	var userId interface{}
	var success bool
	if userId, success = tokenCache.Get(strings.Split(bearerToken, " ")[1]); success != true {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("UserID", userId.(uint))
}

func GetBearerToken(c *gin.Context) {
	var username, password string
	var ok bool
	if username, password, ok = c.Request.BasicAuth(); ok != true {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var user Models.User
	var success bool
	if user, success = Models.GetUser(username); success != true {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Email != username {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token := randstr.String(64)
	userid := user.ID

	duration, err := time.ParseDuration(Configuration.AppConfiguration.TokenLifetime)
	if err != nil {
		duration = 1 * time.Hour
	}
	tokenCache.Set(token, userid, duration)

	c.String(200, "%s", token)
	return

}
