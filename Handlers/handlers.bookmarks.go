package Handlers

import (
	"github.com/canlot/Bookmarkmanager-Server/Helpers"
	"github.com/canlot/Bookmarkmanager-Server/Models"
	"github.com/gin-gonic/gin"
)

func GetBookmarksWithCategoryId(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if bookmarks, success := Models.GetBookmarksWithCategoryID(Helpers.GetUserIDasUint(c), id); success == true {
			c.JSON(200, bookmarks)
		} else {
			Description := "No such bookmarks with this category id: "
			error := Models.JsonError{"No category", Description}
			c.JSON(400, error)
		}

	} else {
		error := Models.JsonError{
			Error:       "Convert Error",
			Description: "Cannot covert id to int, id has a wrong format",
		}
		c.JSON(400, error)
	}
}
func AddBookmarkToCategory(c *gin.Context) {
	idstring := c.Param("category_id")
	var bookmark Models.Bookmark
	if error := c.BindJSON(&bookmark); error == nil {
		if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
			if bookmark, success := Models.AddBookmark(Helpers.GetUserIDasUint(c), id, bookmark); success == true {
				c.JSON(200, bookmark)
			} else {
				error := Models.JsonError{"Database operational failure", ""}
				c.JSON(400, error)
			}
		} else {
			error := Models.JsonError{"Id has not the right format", ""}
			c.JSON(400, error)
		}
	} else {
		error := Models.JsonError{"Request could not be proceed", error.Error()}
		c.JSON(400, error)
	}
}
