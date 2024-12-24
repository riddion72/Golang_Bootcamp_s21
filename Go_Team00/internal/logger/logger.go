package logger

import (
	"log"
	"os"
)

var logFile *os.File

func PrepareLogger(logFileName string) error {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening logfile %s", logFileName)
		return err
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)
	log.Println("Logger initialized")

	return nil
}

func WriteLog(message string) {
	log.Printf("%s\n", message)
}

func CloseLogger() {
	logFile.Close()
}
