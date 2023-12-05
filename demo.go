//+build ignore

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fileStorage := "files"

	year := time.Now().Format("2006")
	month := time.Now().Format("01")

	fileStorage = fileStorage + "/" + year + "/" + month

	if _, err := os.Stat(fileStorage); os.IsNotExist(err) {
		fmt.Println(err)
		os.MkdirAll(fileStorage, os.ModePerm)
	}
	fmt.Println("created")
}
