package utils

import (
	"fmt"
	"strconv"
	"strings"

	"FTP-middle/models"
)

func ParseDirectoryEntry(line string) (entry models.DirectoryEntry, kind int) {
	parts := strings.Fields(line)
	if len(parts) < 7 {
		fmt.Println("Invalid entry:", line)
		return models.DirectoryEntry{}, 0
	}

	kind, _ = strconv.Atoi(parts[1])
	var size int64
	var name string
	if kind == 1 {
		size, _ = strconv.ParseInt(parts[7], 10, 64) // 解析文件大小，假设没有错误处理
		name = strings.Join(parts[8:], " ")
	} else {
		name = strings.Join(parts[7:], " ") // 文件或目录名可能包含空格
	}
	entry = models.DirectoryEntry{
		Permissions: parts[0],
		Owner:       parts[2],
		Group:       parts[3],
		Size:        size,
		Modified:    fmt.Sprintf("%s/%s/%s", parts[6], parts[4], parts[5]),
		Name:        name,
	}
	return
}
