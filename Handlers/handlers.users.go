package Handlers

import (
	"fmt"
	"github.com/canlot/Bookmarkmanager-Server/Helpers"
	"github.com/canlot/Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

func GetUsersForCategory(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if users, success := Models.GetUsersForCategory(id); success == true {
			c.JSON(200, users)
		} else {
			error := Models.JsonError{"No Users", "No Users for this category"}
			c.JSON(400, error)
		}
	} else {
		error := Models.JsonError{"No Users", "No Users for this category"}
		c.JSON(400, error)
	}
}
func AddUsersForCategoryOnce(c *gin.Context) {
	AddUsersForCategory(c, false)
}
func AddUsersForCategoryInherit(c *gin.Context) {
	AddUsersForCategory(c, true)
}

func AddUsersForCategory(c *gin.Context, inherit bool) {
	idstring := c.Param("category_id")
	var users []Models.User
	var id uint
	var success bool
	if id, success = Helpers.ConvertFromStringToUint(&idstring); success == false {
		error := Models.JsonError{"Conversion error", "Id could not be converted"}
		c.JSON(400, error)
		return
	}
	if error := c.BindJSON(&users); error == nil {
		fmt.Println("das ist der Anfang")
		if errstring, success := Models.AddUsersForCategory(Helpers.GetUserIDasUint(c), id, &users, inherit); success == true {
			c.Status(200)
		} else {
			errorjson := Models.JsonError{"ORM error", errstring}
			c.JSON(400, errorjson)
		}
	} else {
		errorjson := Models.JsonError{"Error", error.Error()}
		c.JSON(400, errorjson)
	}
}
func RemoveUsersFromCategory(c *gin.Context) {
	idstring := c.Param("category_id")
	var users []Models.User
	var id uint
	var success bool
	if id, success = Helpers.ConvertFromStringToUint(&idstring); success == false {
		error := Models.JsonError{"Conversion error", "Id could not be converted"}
		c.JSON(400, error)
		return
	}
	if error := c.BindJSON(&users); error == nil {
		if errstring, success := Models.RemoveUsersFromCategory(Helpers.GetUserIDasUint(c), id, &users); success == true {
			c.Status(200)
		} else {
			errorjson := Models.JsonError{"ORM error", errstring}
			c.JSON(400, errorjson)
		}
	} else {
		errorjson := Models.JsonError{"Error", error.Error()}
		c.JSON(400, errorjson)
	}
}

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
