package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ClearSign(c *gin.Context) {
	if Conn != nil {
		Exit(c)
		Conn = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Clear successfully.",
	})
}
