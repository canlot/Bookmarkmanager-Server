package Handlers

import (
	"Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Authenticate(c *gin.Context) {
	authorized := false
	if username, password, ok := c.Request.BasicAuth(); ok == true {
		if user, success := Models.GetUser(username); success == true {
			if user.Email == username {
				if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
					authorized = true
					c.Set("UserID", user.ID)
				} else {
					authorized = false
				}

			} else {
				authorized = false
			}
		}
		if authorized != true {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
