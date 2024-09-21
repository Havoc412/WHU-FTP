package main

import (
	"fmt"
	"net/url"
)

func main() {
	filePath := "/path/to/your/file"
	encodedPath := url.PathEscape(filePath)
	fmt.Println("Encoded file path:", encodedPath)
}
