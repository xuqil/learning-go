package main

import (
	"fmt"
	"os"
)

func main() {
	file := "/root/workplace/goProject/leanring-go/web/testdata/download/myfile.txt"
	fileInfo, _ := os.Stat(file)
	mode := fileInfo.Mode()

	fmt.Println(file, "mode is ", mode)

}
