package handler

import (
	"doctor_doom/helper"
	"doctor_doom/log"
	"github.com/go-co-op/gocron/v2"
	"os"
	"sync"
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

		results := make(chan FileResult, 100) // Buffered channel
		var wg sync.WaitGroup
		sem := make(chan struct{}, 10) // Semaphore to limit concurrency

		log.Debug("Begin to Check Old File")

		// Start the recursive file listing and processing
		wg.Add(1)
		go listFiles(rootPath, minutes, &wg, results, sem)

		// Close the results channel once all processing is done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Handle the results
		for result := range results {
			if result.IsOld {
				log.Debug("File ", result.FilePath, " is older than threshold")
				//deleteFile(result.FilePath)
			} else {
				log.Debug("File ", result.FilePath, " is not older than threshold")
			}
		}
	}
	_, err = s.NewJob(gocron.CronJob("*/1 * * * *", false), gocron.NewTask(deleteTask))
	if err != nil {
		log.Error(err)
	}
	s.Start()
}

func listFiles(dir string, thresholdTime int64, wg *sync.WaitGroup, results chan<- FileResult, sem chan struct{}) {
	defer wg.Done()

	//entries, err := os.ReadDir(dir)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//for _, entry := range entries {
	//	path := filepath.Join(dir, entry.Name())
	//	entry.Info()
	//	if entry.IsDir() {
	//		wg.Add(1)
	//		go listFiles(path, thresholdTime, wg, results, sem)
	//		log.Info("Find out a ", path, " folder")
	//	} else {
	//		wg.Add(1)
	//		go func(filePath string, fileInfo os.DirEntry) {
	//			defer wg.Done()
	//			sem <- struct{}{}        // Acquire token
	//			defer func() { <-sem }() // Release token
	//
	//			info, err := fileInfo.Info()
	//			if err != nil {
	//				log.Error("Error getting file info:", err)
	//				return
	//			}
	//
	//			// Calculate the threshold time
	//			currentTime := time.Now()
	//			threshold := currentTime.Add(-time.Duration(thresholdTime) * time.Minute)
	//
	//			if info.ModTime().Before(threshold) {
	//				// If the file's modification time is before the threshold time, it's considered old
	//				results <- FileResult{FilePath: filePath, IsOld: true}
	//			} else {
	//				// If the file's modification time is after the threshold time, it's considered new
	//				results <- FileResult{FilePath: filePath, IsOld: false}
	//			}
	//		}(path, entry)
	//	}
	//}
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
