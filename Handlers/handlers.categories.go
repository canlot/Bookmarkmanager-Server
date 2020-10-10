package Handlers

import "github.com/gin-gonic/gin"
import "bookmark_Server/Models"

func GetAllCategories(c *gin.Context) {
	categories := Models.GetAllCategories()
	c.JSON(200, categories)

}
