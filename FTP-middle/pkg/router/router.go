package router

import (
	"FTP-middle/pkg/api"
	"FTP-middle/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	v := engine.Group("v1")
	{
		v.GET("hello", api.Hello)
		v.GET("echo/:msg", api.Echo)

		tcp := v.Group("tcp")
		{
			tcp.POST("login", middleware.NotSignInCheck, api.Login)
			tcp.GET("exit", middleware.SignInCheck, api.Exit)

			tcp.GET("clear", api.ClearSign)
		}

		// user := v.Group("user")
		// {
		// 	user.POST("login", api.Login)
		// }

		ftp := v.Group("ftp")
		{
			ftp.GET("list/*path", middleware.SignInCheck, api.List)
			ftp.POST("download", middleware.SignInCheck, api.DownloadFile)
			ftp.POST("upload", middleware.SignInCheck, api.UploadFile)
		}

		ws := v.Group("ws")
		{
			ws.GET("upload", middleware.SignInCheck, api.HandleWSUpload)
			ws.GET("download", middleware.SignInCheck, api.HandleWSDownload)
		}
	}
}
