package middleware

import (
	. "FTP-middle/errcode"
	"FTP-middle/pkg/api"

	"github.com/gin-gonic/gin"
)

func NotSignInCheck(c *gin.Context) {
	if api.Conn != nil {
		c.AbortWithStatusJSON(ErrSignIn, gin.H{
			"err_code": ErrSignIn,
			"err_msg":  ErrMsg[ErrSignIn],
		})
		return
	}
}

func SignInCheck(c *gin.Context) {
	if api.Conn == nil {
		c.AbortWithStatusJSON(ErrNotSignIn, gin.H{
			"err_code": ErrNotSignIn,
			"err_msg":  ErrMsg[ErrNotSignIn],
		})
	}
}
