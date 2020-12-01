package Handlers

import (
	"github.com/canlot/Bookmarkmanager-Server/Helpers"
	"github.com/canlot/Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

func GetUsersForCategoryFull(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if users, success := Models.GetUsersForCategoryFull(id); success == true {
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

func GetUsersForCategoryInherit(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if users, success := Models.GetUsersForCategoryInherit(id); success == true {
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

func AddUsersForCategory(c *gin.Context) {
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
		if errstring, success := Models.AddUsersForCategory(Helpers.GetUserIDasUint(c), id, &users); success == true {
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
