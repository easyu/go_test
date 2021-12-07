package main

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger

func init() {
	format := time.Now().Format("2006-01-02")
	logFile, err := os.OpenFile(format+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("openFile failed,err:%v\n", err)
		return
	}
	logger = log.New(os.Stdout, "<MyServer_1>", log.Lshortfile|log.Ldate|log.Ltime)
	logger.SetOutput(logFile)
	logger.Println("1.--------日志------")
	//log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
func main() {
	logger.Fatal("err")
}
