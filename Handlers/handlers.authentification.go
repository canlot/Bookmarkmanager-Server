package Handlers

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate(c *gin.Context) {
	authorized := false
	if username, password, ok := c.Request.BasicAuth(); ok == true {
		randomBytes := []byte(Helpers.GetRandomString32Lenght())
		passwordBytes := []byte(password)
		all := append(passwordBytes, randomBytes...)
		if user, success := Models.GetUser(username, true); success == true {
			if user.Name == username && user.Password == password {
				authorized = true
				c.Set("UserID", user.ID)
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
