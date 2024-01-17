package main

import (
	"doctor_doom/handler"
	"doctor_doom/log"
	"github.com/labstack/echo/v4"
	"os"
)

func init() {
	os.Setenv("APP_NAME", "XXXX")
	logger := log.InitLogger(false)
	os.Setenv("LOG_LEVEL", "DEBUG")
	logger.SetLevel(log.GetLogLevel("LOG_LEVEL"))
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	deleteFileHandler := handler.DeleteFileHandler{}

	e := echo.New()
	deleteFileHandler.HandlerDeleteFile()
	e.Logger.Fatal(e.Start(":1323"))
}
