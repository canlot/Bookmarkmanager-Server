package Helpers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ConvertFromStringToUint(s *string) (number uint, success bool) {
	if n, err := strconv.ParseUint(*s, 10, 32); err == nil {
		return uint(n), true
	} else {
		return 0, false
	}
}

func GetUserIDasUint(c *gin.Context) uint {
	result, _ := c.Get("UserID")
	id := result.(uint)
	return id
}
