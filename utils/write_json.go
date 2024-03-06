package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func WriteToJson(data []byte) {
	currentPath, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		currentPath += "\\"
	} else if runtime.GOOS == "linux" {
		currentPath += "/"
	}
	fileName := currentPath + time.Now().Format("2006-01-02_15-04-05-1") + fmt.Sprintf("_record.json")
	var file, err = os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("Failed to create file %s, %v", fileName, err)
	}
	file.Write(data)
	log.Println("==============Export to json file success", fileName)
}
