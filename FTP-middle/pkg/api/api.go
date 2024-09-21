package api

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"FTP-middle/consts"
)

var Conn net.Conn

type (
	ErrorRes struct {
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
	}
)

type (
	E ErrorRes
)

const (
	TargetIP = "localhost"
	// TargetIP = "26.177.188.34" // test 测试 VPN；successful
)

func Hello(c *gin.Context) { // info 测试用
	c.JSON(http.StatusOK, gin.H{
		"msg": "Hello World!",
	})
}

func Echo(c *gin.Context) { // update
	conn, err := net.Dial("tcp", "localhost:8000") // info link to Server
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close() // tip 确保在退出前关闭连接

	msg := c.Param("msg")
	fmt.Fprintf(conn, consts.ECHO+msg)                       // send
	res, err := bufio.NewReader(conn).ReadString(consts.EOF) // recive
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": res[:len(res)-1],
	})
}

func sendCommand_getMessage(cmd string) (res string) {
	fmt.Fprint(Conn, cmd)
	res, err := bufio.NewReader(Conn).ReadString(consts.EOF)
	if err != nil {
		fmt.Println("Error reading response: ", err.Error())
		return // todo 补充报错反馈
	}
	return
}

func getBackMessage() (res string) {
	res, err := bufio.NewReader(Conn).ReadString(consts.EOF)
	if err != nil {
		fmt.Println("Error reading response: ", err.Error())
		return // todo 补充报错反馈
	}
	return
}
