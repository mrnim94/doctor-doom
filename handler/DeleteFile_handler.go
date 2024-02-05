package handler

import (
	"doctor_doom/helper"
	"doctor_doom/log"
	"github.com/go-co-op/gocron/v2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DeleteFileHandler struct {
}

type FileResult struct {
	FilePath string
	IsOld    bool
}

func (dl *DeleteFileHandler) HandlerDeleteFile() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Error(err)
	}

	deleteTask := func() {
		rootPath := helper.GetEnvOrDefault("DOOM_PATH", "test_delete")

		minutes := helper.DurationToMinutes(helper.GetEnvOrDefault("RULE_AGE", "1h"))
		log.Debug("Begin to Check Old File")

		resultsChan := make(chan string, 100) // Buffered channel for file paths to delete
		doneChan := make(chan bool)

		go func() {
			for filePath := range resultsChan {
				deleteFile(filePath)
			}
			doneChan <- true
		}()

		var wg sync.WaitGroup
		wg.Add(1)
		go listFiles(&wg, rootPath, minutes, resultsChan)

		wg.Wait()
		close(resultsChan) // Close the results channel to signal the deletion goroutine to finish
		<-doneChan         // Wait for the deletion goroutine to signal it's done
		log.Debug("Completed file processing.")
	}
	_, err = s.NewJob(gocron.CronJob(helper.GetEnvOrDefault("CIRCLE", "*/1 * * * *"), false), gocron.NewTask(deleteTask))
	if err != nil {
		log.Error(err)
	}
	s.Start()
}

func listFiles(wg *sync.WaitGroup, dir string, thresholdTime int64, resultsChan chan<- string) {
	defer wg.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Error("Error reading directory: ", err)
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			wg.Add(1)
			go listFiles(wg, path, thresholdTime, resultsChan)
		} else {
			info, err := entry.Info()
			if err != nil {
				log.Error("Error getting file info: ", err)
				continue
			}
			currentTime := time.Now()
			threshold := currentTime.Add(-time.Duration(thresholdTime) * time.Minute)
			if info.ModTime().Before(threshold) {
				resultsChan <- path
			}
		}
	}
}

func deleteFile(filePath string) {
	// Check if the file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Debug("File does not exist: ", filePath)
		} else {
			log.Error("Error checking file existence: ", err)
		}
		return
	}

	// Attempt to delete the file
	err := os.Remove(filePath)
	if err != nil {
		log.Error("Error deleting the file: ", err)
		return
	}

	log.Debug("File deleted successfully: ", filePath)
}
