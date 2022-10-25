package main

import (
	"fmt"
)

var stor Storage

func main() {
	fmt.Println("Hello")

	db := SqlLiteDb{DBFilePath: "./data/test.db"}
	stor = &db
}
