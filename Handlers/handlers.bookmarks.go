package Handlers

import (
	"Bookmarkmanager-Server/Configuration"
	"Bookmarkmanager-Server/Helpers"
	"Bookmarkmanager-Server/Models"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
)

func GetBookmarksWithCategoryId(c *gin.Context) {
	idstring := c.Param("category_id")
	if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
		if bookmarks, success := Models.GetBookmarksWithCategoryID(Helpers.GetUserIdAsUint(c), id); success == true {
			c.JSON(200, bookmarks)
			return
		} else {
			Description := "No such bookmarks with this category id: "
			error := Models.JsonError{"No category", Description}
			c.JSON(400, error)
			return
		}

	} else {
		error := Models.JsonError{
			Error:       "Convert Error",
			Description: "Cannot covert id to int, id has a wrong format",
		}
		c.JSON(400, error)
		return
	}
}
func SearchBookmarks(c *gin.Context) {
	searchText := c.Param("search_text")
	if len(searchText) < 2 {
		c.JSON(400, errors.New("At least 3 characters"))
		return
	}
	bookmarks, err := Models.SearchBookmarks(Helpers.GetUserIdAsUint(c), searchText)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, bookmarks)
}
func AddBookmarkToCategory(c *gin.Context) {
	idstring := c.Param("category_id")
	var bookmark Models.Bookmark
	if error := c.BindJSON(&bookmark); error == nil {
		if id, success := Helpers.ConvertFromStringToUint(&idstring); success == true {
			if bookmark, success := Models.AddBookmark(Helpers.GetUserIdAsUint(c), id, bookmark); success == true {
				c.JSON(200, bookmark)
				return
			} else {
				error := Models.JsonError{"Database operational failure", ""}
				c.JSON(400, error)
				return
			}
		} else {
			error := Models.JsonError{"Id has not the right format", ""}
			c.JSON(400, error)
			return
		}
	} else {
		error := Models.JsonError{"Request could not be proceed", error.Error()}
		c.JSON(400, error)
		return
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
		return
	}
	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	var bookmark Models.Bookmark
	if err := c.BindJSON(&bookmark); err != nil {
		c.JSON(400, "Request could not be proceed")
		return
	}

	if bookmarkid != bookmark.ID {
		c.JSON(400, errors.New("Different ids"))
		return
	}

	err := Models.EditBookmark(Helpers.GetUserIdAsUint(c), categoryid, bookmark)

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.Status(200)

}
func DeleteBookmarkWithBookmarkId(c *gin.Context) {
	categoryidstring := c.Param("category_id")
	bookmarkidstring := c.Param("bookmark_id")

	var categoryid uint
	var bookmarkid uint
	var success bool

	if categoryid, success = Helpers.ConvertFromStringToUint(&categoryidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}
	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	err := Models.DeleteBookmark(Helpers.GetUserIdAsUint(c), categoryid, bookmarkid)

	if err != nil {
		c.JSON(400, err)
	}
	c.Status(200)
}
func MoveBookmarkWithBookmarkId(c *gin.Context) {
	categoryidstring := c.Param("category_id")
	bookmarkidstring := c.Param("bookmark_id")
	categorydestinationidstring := c.Param("category_destination_id")

	var sourceid uint
	var destinationid uint
	var bookmarkid uint
	var success bool

	if sourceid, success = Helpers.ConvertFromStringToUint(&categoryidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}
	if destinationid, success = Helpers.ConvertFromStringToUint(&categorydestinationidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}
	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}
	err := Models.MoveBookmarkToCategory(Helpers.GetUserIdAsUint(c), bookmarkid, sourceid, destinationid)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.Status(200)
}
func GetIconForBookmark(c *gin.Context) {
	bookmarkidstring := c.Param("bookmark_id")
	categoryidstring := c.Param("category_id")
	var bookmarkid, categoryid, userid uint
	var success bool

	userid = Helpers.GetUserIdAsUint(c)

	if categoryid, success = Helpers.ConvertFromStringToUint(&categoryidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	bookmark, err := Models.GetBookmark(userid, categoryid, bookmarkid)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if bookmark.IconHash == "" {
		c.JSON(400, errors.New("No icon for this bookmark"))
		return
	}
	fileName := bookmark.IconHash + ".png"
	iconPath := path.Join(Configuration.AppConfiguration.IconFolderPath, fileName)

	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		c.JSON(400, errors.New("No icon for this bookmark"))
		return
	}
	rawfile, err := os.ReadFile(iconPath)
	if err != nil {
		log.Print(err)
		c.JSON(400, err)
		return
	}
	c.Data(200, "image/png", rawfile)
	//c.FileAttachment(iconPath, fileName)
	return

}
func UploadIconToBookmark(c *gin.Context) {
	bookmarkidstring := c.Param("bookmark_id")
	categoryidstring := c.Param("category_id")
	var bookmarkid, categoryid, userid uint
	var success bool

	userid = Helpers.GetUserIdAsUint(c)

	if categoryid, success = Helpers.ConvertFromStringToUint(&categoryidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	if bookmarkid, success = Helpers.ConvertFromStringToUint(&bookmarkidstring); success != true {
		c.JSON(400, errors.New("Could not convert string"))
		return
	}

	bookmark, err := Models.GetBookmark(userid, categoryid, bookmarkid)
	if err != nil {
		c.JSON(400, err)
		return
	}

	if c.ContentType() != "image/png" {
		c.JSON(400, errors.New("Wrong filetype"))
		return
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(400, err)
		return
	}
	hash, err := Helpers.CumputeHashFromBytes(raw)
	if err != nil {
		c.JSON(400, err)
		return
	}

	fileName := hash + ".png"

	iconPath := path.Join(Configuration.AppConfiguration.IconFolderPath, fileName)

	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		file, err := os.Create(iconPath)
		if err != nil {
			c.JSON(400, err)
			return
		}
		defer file.Close()
		_, err = file.Write(raw)
		if err != nil {
			c.JSON(400, err)
			return
		}
	}

	bookmark.IconName = fileName

	err = Models.EditBookmark(userid, categoryid, bookmark)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.Status(200)
}

func UploadTest(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(400, err)
		return
	}
	hash, err := Helpers.CumputeHashFromBytes(raw)
	if err != nil {
		c.JSON(400, err)
		return
	}

	fileName := hash + ".png"

	iconPath := path.Join(Configuration.AppConfiguration.IconFolderPath, fileName)
	file, err := os.Create(iconPath)
	defer file.Close()
	file.Write(raw)
	c.Status(200)
	return
}
