package Handlers

import (
	"Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	if users, success := Models.GetAllUsers(); success == true {
		c.JSON(200, users)
	} else {
		c.Status(400)
	}
}
func GetCurrentUser(c *gin.Context) {
	if username, _, ok := c.Request.BasicAuth(); ok == true {
		if user, success := Models.GetUser(username, false); success == true {
			c.JSON(200, user)
		} else {
			c.Status(400)
		}
	}
}
