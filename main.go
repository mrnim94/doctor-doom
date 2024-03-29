package main

import (
	"doctor_doom/handler"
	"doctor_doom/log"
	"github.com/labstack/echo/v4"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func init() {
	os.Setenv("APP_NAME", "doctor-doom")
	logger := log.InitLogger(false)
	// Check if KUBERNETES_SERVICE_HOST is set
	if _, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST"); !exists {
		// If not in Kubernetes, set LOG_LEVEL to DEBUG
		os.Setenv("LOG_LEVEL", "DEBUG")
	}
	logger.SetLevel(log.GetLogLevel("LOG_LEVEL"))
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}

func main() {

	deleteFileHandler := handler.DeleteFileHandler{}

	e := echo.New()

	deleteFileHandler.HandlerDeleteFile()
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			e.Logger.Fatal("pprof server failed to start:", err)
		}
	}()
	e.Logger.Fatal(e.Start(":1994"))
}
