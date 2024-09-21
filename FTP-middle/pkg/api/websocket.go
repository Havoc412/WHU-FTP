package api

import (
	"FTP-middle/consts"
	"FTP-middle/models"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "FTP-middle/errcode"
	. "FTP-middle/redis"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // info 在生产环境中可能需要更安全的检查
	},
}

func HandleWSUpload(c *gin.Context) {
	// 1. get arguments
	sessionId := c.Query("sessionId")
	val, err := Rdb.Get(Ctx, sessionId).Result()

	if err != nil {
		log.Println("Error retrieving seesion data:", err)
		c.JSON(ErrGetRides, E{
			ErrCode: ErrGetRides,
			ErrMsg:  ErrMsg[ErrGetRides],
		})
		return
	}

	var uploadFile models.UploadFile
	if err := json.Unmarshal([]byte(val), &uploadFile); err != nil {
		log.Println("Error unmarshalling session data", err)
		c.JSON(ErrGetRides, E{
			ErrCode: ErrGetRides,
			ErrMsg:  ErrMsg[ErrGetRides], // 先直接用相同的报错
		})
		return
	}
	// log.Println(uploadFile.TargetPath, uploadFile.LocalFilePath)

	// info 2. prepare websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer ws.Close()

	// info 3. open local file
	file, err := os.Open(uploadFile.LocalFilePath)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		c.JSON(ErrOpenLocalFile, E{
			ErrCode: ErrOpenLocalFile,
			ErrMsg:  ErrMsg[ErrOpenLocalFile],
		})
		return
	}
	defer file.Close()

	// 3. send cmd to FTP
	fileName := filepath.Base(uploadFile.LocalFilePath)
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	fmt.Fprintf(Conn, consts.UPLOAD+"%s %d %s", uploadFile.TargetPath, fileSize, fileName)

	// 4. 获取新的端口号 并建立连接。
	res := getBackMessage()
	parts := strings.Split(strings.TrimSpace(res), " ")
	if parts[0] != "UPLOAD_PORT" {
		fmt.Println("Error reading response:", err)
		c.JSON(ErrGetNewLink, E{
			ErrCode: ErrGetNewLink,
			ErrMsg:  ErrMsg[ErrGetNewLink],
		})
		return
	}
	newPort := parts[1]
	fmt.Println("New transfer port:", newPort)
	transferConn, err := net.Dial("tcp", TargetIP+":"+newPort)
	if err != nil {
		fmt.Println("Error connecting to transfer port:", err)
		return
	}
	defer transferConn.Close()

	// info CORE 这里开一个 协程 监听反向的 ws；
	controlChan := make(chan bool) // tip 全局监听的控制信号
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading ws message:", err)
				return
			}

			var cmd models.VueCmd
			if err := json.Unmarshal(message, &cmd); err != nil {
				log.Println("Error unmarshalign message:", err)
				continue
			}

			switch cmd.Command {
			case "pause":
				controlChan <- false
			case "resume":
				controlChan <- true
			default:
				log.Println("Unknown command:", cmd.Command)
			}
		}
	}()

	// 4. send data
	buffer := make([]byte, 1024*1024) // 1MB buffer
	totalBytesSent := int64(0)
	paused := false
	for {
		if paused {
			if <-controlChan {
				_, writeErr := transferConn.Write([]byte(consts.FTP_RESUME))
				log.Println("RESUME sended")
				if writeErr != nil {
					fmt.Println("Error resume task:", writeErr)
					c.JSON(ErrTaskResume, E{
						ErrCode: ErrTaskResume,
						ErrMsg:  ErrMsg[ErrTaskResume],
					})
					return
				}
				paused = false
			}
		} else {
			select {
			case ctrl := <-controlChan:
				if !ctrl {
					paused = true
					_, writeErr := transferConn.Write([]byte(consts.FTP_PAUSE))
					if writeErr != nil {
						fmt.Println("Error pause task:", writeErr)
						c.JSON(ErrTaskPause, E{
							ErrCode: ErrTaskPause,
							ErrMsg:  ErrMsg[ErrTaskPause],
						})
						return
					}
					continue
				}
			default:
				bytesRead, readErr := file.Read(buffer)
				if readErr != nil {
					if readErr == io.EOF {
						log.Println("read over")
						goto FinishUpload // tip return、break 好像都达不到效果...
					}
					fmt.Println("Error reading file:", readErr)

					// todo 错误断点处理
					return
				}
				// 发送文件数据，加入前缀 "DATA:"
				dataWithPrefix := fmt.Sprintf(consts.FTP_DATA_PREFIX+"%s", buffer[:bytesRead])
				_, writeErr := transferConn.Write([]byte(dataWithPrefix))
				if writeErr != nil {
					fmt.Println("Error sending file:", writeErr)
					// todo 处理 中断重连
					return
				}

				totalBytesSent += int64(bytesRead)
				// ws 保持更新
				ws.WriteJSON(models.WsState{
					State:    1,
					SentByte: totalBytesSent,
				})
			}
		}
	}
FinishUpload: // 这里之后同时也顺利执行了 ws.Close() 的样子。
	log.Println("File sent successfully.")

	res = getBackMessage()
	// res := "123\n" // test for no FTP, because getBackMessage lead to blockking
	// 使用 WebSocket 发送数据
	msg := res[:len(res)-1] // 根据需要处理消息
	if err := ws.WriteJSON(models.WsState{
		State:   2,
		Message: msg,
	}); err != nil {
		log.Println("Failed to send message over websocket:", err)
	}
}

func HandleWSDownload(c *gin.Context) {
	// 0.
	sessionId := c.Query("sessionId")
	val, err := Rdb.Get(Ctx, sessionId).Result()
	if err != nil {
		log.Println("Error retrieving seesion data:", err)
		c.JSON(ErrGetRides, E{
			ErrCode: ErrGetRides,
			ErrMsg:  ErrMsg[ErrGetRides],
		})
		return
	}
	var downloadFile models.DownloadFile
	if err := json.Unmarshal([]byte(val), &downloadFile); err != nil {
		log.Println("Error unmarshalling session data", err)
		c.JSON(ErrGetRides, E{
			ErrCode: ErrGetRides,
			ErrMsg:  ErrMsg[ErrGetRides], // 先直接用相同的报错
		})
		return
	}

	// 1. open Websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer ws.Close()

	// 2. open the file
	file, err := os.Create(downloadFile.SavePath)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		c.JSON(ErrCreateFile, E{
			ErrCode: ErrCreateFile,
			ErrMsg:  ErrMsg[ErrCreateFile],
		})
		return
	}
	defer file.Close()
	log.Println(downloadFile.TargetPath, downloadFile.SavePath) // test

	// 3. send cmd && get new port
	fmt.Fprintf(Conn, consts.DOWNLOAD+downloadFile.TargetPath) // send

	// 4. 获取新的端口号 并建立连接。
	res := getBackMessage()
	parts := strings.Split(strings.TrimSpace(res), " ")
	if parts[0] != "DOWNLOAD_PORT" {
		fmt.Println("Error reading response:", err)
		c.JSON(ErrGetNewLink, E{
			ErrCode: ErrGetNewLink,
			ErrMsg:  ErrMsg[ErrGetNewLink],
		})
		return
	}
	newPort := parts[1]
	fmt.Println("New transfer port:", newPort)
	transferConn, err := net.Dial("tcp", TargetIP+":"+newPort) // todo 放到全局。
	if err != nil {
		fmt.Println("Error connecting to transfer port:", err)
		return
	}
	defer transferConn.Close()

	// 5. CORE TRANSFORM: first step
	fileSizeBytes := make([]byte, 8)
	_, err = io.ReadFull(transferConn, fileSizeBytes)
	if err != nil {
		fmt.Println("Error reading file size: ", err)
		c.JSON(ErrDownloadFileLength, E{
			ErrCode: ErrDownloadFileLength,
			ErrMsg:  ErrMsg[ErrDownloadFileLength],
		})
		return
	}
	fileSize := int64(binary.LittleEndian.Uint64(fileSizeBytes))

	// info second step
	// 6-middle. 建立完整的 websocket 连接
	controlChan := make(chan bool) // tip 全局监听的控制信号
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading ws message:", err)
				return
			}

			var cmd models.VueCmd
			if err := json.Unmarshal(message, &cmd); err != nil {
				log.Println("Error unmarshalign message:", err)
				continue
			}

			switch cmd.Command {
			case "pause":
				controlChan <- false
			case "resume":
				controlChan <- true
			default:
				log.Println("Unknown command:", cmd.Command)
			}
		}
	}()

	buffer := make([]byte, 1024*1024)
	var receivedBytes int64
	paused := false
	for receivedBytes < fileSize {
		// log.Println("running", receivedBytes, paused)
		if paused {
			if <-controlChan {
				_, writeErr := transferConn.Write([]byte(consts.FTP_RESUME)) //
				log.Println("RESUME sended")
				if writeErr != nil {
					fmt.Println("Error resume task:", writeErr)
					c.JSON(ErrTaskResume, E{
						ErrCode: ErrTaskResume,
						ErrMsg:  ErrMsg[ErrTaskResume],
					})
					return
				}
				paused = false
			}
		} else {
			select {
			case ctrl := <-controlChan:
				if !ctrl {
					paused = true
					_, writeErr := transferConn.Write([]byte(consts.FTP_PAUSE))
					if writeErr != nil {
						fmt.Println("Error pause task:", writeErr)
						c.JSON(ErrTaskPause, E{
							ErrCode: ErrTaskPause,
							ErrMsg:  ErrMsg[ErrTaskPause],
						})
						return
					}
					continue
				}
			default:
				bytesRead, err := transferConn.Read(buffer)
				if err != nil {
					if err == io.EOF {
						goto FinishDownloadFinish
					}
					fmt.Errorf("failed to read from connection: %w", err)
					// todo 错误处理
					return
				}

				if _, err := file.Write(buffer[:bytesRead]); err != nil {
					fmt.Errorf("failed to write to file: %w", err)
					// todo
					return
				}
				receivedBytes += int64(bytesRead)
				// ws 保持更新
				ws.WriteJSON(models.WsState{
					State:    1,
					SentByte: receivedBytes,
				})
			}
		}
	}

FinishDownloadFinish:
	res = getBackMessage()

	msg := res[:len(res)-1] // 根据需要处理消息
	if err := ws.WriteJSON(models.WsState{
		State:   2,
		Message: msg,
	}); err != nil {
		log.Println("Failed to send message over websocket:", err)
	}
}
