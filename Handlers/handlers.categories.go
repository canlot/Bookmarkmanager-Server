package Handlers

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"errors"
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
func SearchCategories(c *gin.Context) {
	searchText := c.Param("search_text")
	if len(searchText) < 2 {
		c.JSON(400, errors.New("At least 3 characters"))
	}
	categories, err := Models.SearchCategories(Helpers.GetUserIdAsUint(c), searchText)
	if err != nil {
		c.JSON(400, err)
	}
	c.JSON(200, categories)
}

func AddCategory(c *gin.Context) {
	var category Models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, err)
		return
	}
	if returnCategory, err := Models.AddCategory(Helpers.GetUserIdAsUint(c), category); err != nil {
		c.JSON(400, err)
		return
	} else {
		c.JSON(200, returnCategory)
		return
	}

}
func EditCategory(c *gin.Context) {
	var category Models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, err)
		return
	}
	if returnCategory, err := Models.EditCategory(Helpers.GetUserIdAsUint(c), category); err != nil {
		c.JSON(400, err)
		return
	} else {
		c.JSON(200, returnCategory)
	}
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

func EditUsersForCategory(c *gin.Context) {
	idstring := c.Param("category_id")
	var users []Models.User
	var id uint
	var success bool
	if id, success = Helpers.ConvertFromStringToUint(&idstring); success == false {
		c.JSON(400, errors.New("Id could not be converted"))
		return
	}
	if err := c.BindJSON(&users); err != nil {
		c.JSON(400, err)
		return
	}
	if err := Models.EditUsersForCategory(Helpers.GetUserIdAsUint(c), id, &users); err != nil {
		c.JSON(400, err)
		return
	}
	c.Status(200)
	return
}
