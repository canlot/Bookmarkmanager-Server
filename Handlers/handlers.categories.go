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
			return
		}
	}

	if categories, err := Models.GetCategories(Helpers.GetUserIdAsUint(c), categoryId); err == nil {
		c.JSON(200, categories)
		return
	} else {
		c.JSON(400, err)
		return
	}

}

func AddCategory(c *gin.Context) {
	var category Models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, err)
		return
	}
	if err := Models.AddCategory(Helpers.GetUserIdAsUint(c), category); err != nil {
		c.JSON(400, err)
		return
	} else {
		c.Status(200)
		return
	}

}
func EditCategory(c *gin.Context) {
	var category Models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, err)
		return
	}
	if err := Models.EditCategory(Helpers.GetUserIdAsUint(c), category); err != nil {
		c.JSON(400, err)
		return
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
	if Models.DeleteCategory(categoryId, Helpers.GetUserIdAsUint(c)) != true {
		c.Status(400)
		return
	}
	c.Status(200)
}
