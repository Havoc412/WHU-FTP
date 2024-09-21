package api

import (
	"FTP-middle/consts"
	"FTP-middle/models"
	"FTP-middle/pkg/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	. "FTP-middle/redis"
)

func List(c *gin.Context) {
	encodedPath := c.Param("path")
	path, err := url.PathUnescape(encodedPath) // info url -> path
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}
	fmt.Println(encodedPath, path)

	var dirList []models.DirectoryEntry
	var fileList []models.DirectoryEntry

	// res := sendCommand(conn, consts.LIST+path)
	fmt.Fprintf(Conn, consts.LIST+path) // send

	// info scan data line by line
	scanner := bufio.NewScanner(Conn)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		if line == consts.STOP {
			fmt.Println("list stop.")
			break
		}
		entry, kind := utils.ParseDirectoryEntry(line)
		// fmt.Println(entry)
		if kind == 1 {
			fileList = append(fileList, entry)
		} else if kind == 2 {
			dirList = append(dirList, entry)
		} // kind = 0 过滤
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "list over.",
		"dir_list":  dirList,
		"file_list": fileList,
	})
}

func DownloadFile(c *gin.Context) {
	var downloadFile models.DownloadFile
	if err := c.ShouldBind(&downloadFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// todo 这里同时可以处理其他的一些需求设计
	sessionID := utils.GenerateSessionID()

	json, _ := json.Marshal(downloadFile)
	Rdb.Set(Ctx, sessionID, json, time.Hour*1)

	c.JSON(http.StatusOK, gin.H{
		"msg":       "Download ready, connect via WebSocket",
		"sessionId": sessionID,
	})
}

func UploadFile(c *gin.Context) {
	var uploadFile models.UploadFile
	if err := c.ShouldBind(&uploadFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// todo 这里同时可以处理其他的一些需求设计
	sessionID := utils.GenerateSessionID()

	json, _ := json.Marshal(uploadFile)
	Rdb.Set(Ctx, sessionID, json, time.Hour*1)

	c.JSON(http.StatusOK, gin.H{
		"msg":       "Upload ready, connect via WebSocket",
		"sessionId": sessionID,
	})
}
