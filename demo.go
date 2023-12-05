// +build ignore

package main

import (
	"fmt"
	"os"
)

func main() {
	// 创建文件夹
	/*fileStorage := "files"
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	fileStorage = fileStorage + "/" + year + "/" + month
	if _, err := os.Stat(fileStorage); os.IsNotExist(err) {
		fmt.Println(err)
		os.MkdirAll(fileStorage, os.ModePerm)
	}
	fmt.Println("created")*/

	// 获取文件列表
	/*files, err := filepath.Glob(filepath.Join(config.FileStorage, "*"))
	if err != nil {
		panic(err)
	}

	var links []string
	for _, file := range files {
		links = append(links, strings.TrimPrefix(file, config.FileStorage+"/"))
	}
	fmt.Println(links)*/

	// 删除文件
	filename := "files\\2023\\12\\068e4a7b6498fc7aefc6f3a6ef34e45e.jpg"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("文件不存在")
	} else {
		fmt.Println("文件存在")
	}

	err := os.Remove(filename)
	if err != nil {
		fmt.Println("删除失败:", err)
	}
	fmt.Println("删除成功")
}
