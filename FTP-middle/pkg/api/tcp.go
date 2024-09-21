package api

import (
	"FTP-middle/consts"
	"FTP-middle/models"
	"FTP-middle/pkg/utils"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	. "FTP-middle/errcode"
)

func Login(c *gin.Context) {
	var login models.Login
	// username, password := c.Query("username"), c.Query("password")
	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Print(login.Username, login.Password) // todo 之后实现具体的鉴权。

	var err error
	Conn, err = net.Dial("tcp", TargetIP+":8000")
	if err != nil {
		fmt.Println("Error connecting: ", err.Error())
		c.JSON(ErrServer, E{
			ErrCode: ErrServer,
			ErrMsg:  ErrMsg[ErrServer],
		})
		return
	}
	// info encrypt string
	password, err := utils.EncryptString_AES256(login.Password, consts.AES256_KEY, consts.AES256_IV)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	res := sendCommand_getMessage(consts.LOGIN + login.Username + " " + password) // todo 密码需要加密

	c.JSON(http.StatusOK, gin.H{
		"msg": res[:len(res)-1],
	})
}

func Exit(c *gin.Context) {
	if Conn != nil {
		res := sendCommand_getMessage(consts.EXIT)
		fmt.Println("[EXIT] ", res)

		err := Conn.Close()
		if err != nil {
			log.Println("Failed to close the connection:", err)
		} else {
			log.Println("Connection closed successfully")
		}
		Conn = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "退出成功。",
	})
}

func getNewLink(port string) {

}
