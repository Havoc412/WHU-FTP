package main

import (
	"FTP-middle/pkg/router"
	// _ "FTP-midlle/errcode" // question 单纯先初始化

	_ "FTP-middle/redis"

	"github.com/gin-gonic/gin"
)

// info CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	engine := gin.Default()

	engine.Use(CORSMiddleware()) // test 应用 CORS 中间件

	router.Register(engine)
	port := "8080"
	engine.Run(":" + port)
}
