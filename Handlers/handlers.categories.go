package Handlers

import "github.com/gin-gonic/gin"
import "github.com/canlot/Bookmarkmanager-Server/Models"

func GetAllCategories(c *gin.Context) {
	categories := Models.GetAllCategories()
	c.JSON(200, categories)

}
