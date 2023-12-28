package Handlers

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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
		if user, success := Models.GetUser(username); success == true {
			c.JSON(200, user)
		} else {
			c.Status(400)
		}
	}
}
func AddUser(c *gin.Context) {
	password := c.Param("password")
	var err error
	var user Models.User
	if err = c.BindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	if user, err = Models.AddUser(Helpers.GetUserIdAsUint(c), user, password); err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
func SearchUsers(c *gin.Context) {
	seachString := c.Param("search_text")
	if len(seachString) < 3 {
		c.JSON(400, errors.New("search text was too short"))
	}
	var err error
	var users []Models.User
	if users, err = Models.SearchUsers(seachString); err != nil {
		c.JSON(400, err)
	}
	c.JSON(200, users)
}
