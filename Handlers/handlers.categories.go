package Handlers

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	stringid := c.Param("category_id")
	var categoryId uint
	if stringid == "" {
		categoryId = 0
	} else {
		var success bool
		if categoryId, success = Helpers.ConvertFromStringToUint(&stringid); success != true {
			error := Models.JsonError{"Wrong category", "Could not convert category"}
			c.JSON(400, error)
		}
	}
	if categories, success := Models.GetCategories(Helpers.GetUserIDasUint(c), categoryId); success == true {
		c.JSON(200, categories)
	} else {
		error := Models.JsonError{"Could not fetch", "kdjsl"}
		c.JSON(400, error)
	}

}

func AddCategory(c *gin.Context) {
	var category Models.Category
	if error := c.BindJSON(&category); error == nil {
		if category, success := Models.AddCategory(Helpers.GetUserIDasUint(c), category); success == true {
			c.JSON(200, category)
		} else {
			c.Status(400)
		}

	} else {
		error := Models.JsonError{
			Error:       "Request could not be proceed",
			Description: error.Error(),
		}
		c.JSON(400, error)
	}

}
func EditCategory(c *gin.Context) {
	var category Models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, err)
		return
	}
	if err := Models.EditCategory(Helpers.GetUserIDasUint(c), category); err != nil {
		c.JSON(400, err)
	}
	c.Status(200)
}

func DeleteCategory(c *gin.Context) {
	stringid := c.Param("category_id")
	var categoryId uint
	if stringid == "" {
		c.Status(400)
		return
	}
	var success bool
	if categoryId, success = Helpers.ConvertFromStringToUint(&stringid); success != true {
		c.Status(400)
		return
	}
	if Models.DeleteCategory(categoryId, Helpers.GetUserIDasUint(c)) != true {
		c.Status(400)
		return
	}
	c.Status(200)
}
