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
		return
	} else {
		c.Status(400)
		return
	}
}
func GetCurrentUser(c *gin.Context) {
	if username, _, ok := c.Request.BasicAuth(); ok == true {
		if user, success := Models.GetUser(username); success == true {
			c.JSON(200, user)
			return
		} else {
			c.Status(400)
			return
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
func EditUser(c *gin.Context) {
	userIdString := c.Param("id")
	var userId uint
	var success bool
	if userId, success = Helpers.ConvertFromStringToUint(&userIdString); success != true {
		c.JSON(400, errors.New("couldn't convert user id"))
		return
	}
	password := c.Param("password")
	var user Models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, errors.New("couldn't parse user object"))
		return
	}
	if userId != user.ID {
		c.JSON(400, errors.New("user ids not matching"))
		return
	}
	if err := Models.EditUser(Helpers.GetUserIdAsUint(c), user, password); err != nil {
		c.JSON(400, err)
		return
	}
	c.Status(200)

}
func DeleteUser(c *gin.Context) {
	userIdString := c.Param("id")
	var userId uint
	var success bool
	if userId, success = Helpers.ConvertFromStringToUint(&userIdString); success != true {
		c.JSON(400, errors.New("couldn't convert user id"))
		return
	}
	if err := Models.DeleteUser(Helpers.GetUserIdAsUint(c), userId); err != nil {
		c.JSON(400, err)
	}
	c.Status(200)

}
func SearchUsers(c *gin.Context) {
	searchString := c.Param("search_text")
	if len(searchString) < 3 {
		c.JSON(400, errors.New("search text was too short"))
		return
	}
	var err error
	var users []Models.User
	if users, err = Models.SearchUsers(searchString); err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, users)
}
