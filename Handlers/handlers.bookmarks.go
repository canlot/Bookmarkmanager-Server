package Handlers

import (
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetBookmarksWithCategoryId(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if bookmarks, success := Models.GetBookmarksWithCategoryID(Helpers.GetUserIdAsUint(c), id); success == true {
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
			if bookmark, success := Models.AddBookmark(Helpers.GetUserIdAsUint(c), id, bookmark); success == true {
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
func EditBookmarkWithBookmarkId(c *gin.Context) {
	categoryidstring := c.Param("category_id")
	bookmarkidstring := c.Param("bookmark_id")

	var categoryid uint
	var bookmarkid uint
	var success bool

	if categoryid, success = Helpers.ConvertFromStringToUint(&categoryidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
	}
	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
	}

	var bookmark Models.Bookmark
	if error := c.BindJSON(&bookmark); error != nil {
		c.JSON(400, "Request could not be proceed")
	}

	if bookmarkid != bookmark.ID {
		c.JSON(400, errors.New("Different ids"))
	}

}
func DeleteBookmarkWithBookmarkId(c *gin.Context) {
	categoryidstring := c.Param("category_id")
	bookmarkidstring := c.Param("bookmark_id")
}
