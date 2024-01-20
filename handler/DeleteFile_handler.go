package handler

import (
	"doctor_doom/helper"
	"doctor_doom/log"
	"github.com/go-co-op/gocron/v2"
	"os"
	"path/filepath"
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

		// Start the recursive file listing and processing
		var results []FileResult
		results, err = listFiles(rootPath, minutes, results)

		// Close the results channel once all processing is done

		// Handle the results
		for _, result := range results {
			if result.IsOld {
				log.Debug("File ", result.FilePath, " is older than threshold")
				deleteFile(result.FilePath)
			} else {
				log.Debug("File ", result.FilePath, " is not older than threshold")
			}
		}
	}
	_, err = s.NewJob(gocron.CronJob(helper.GetEnvOrDefault("CIRCLE", "*/1 * * * *"), false), gocron.NewTask(deleteTask))
	if err != nil {
		log.Error(err)
	}
	s.Start()
}

func listFiles(dir string, thresholdTime int64, results []FileResult) ([]FileResult, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			// It's a directory, recurse into it
			var err error
			results, err = listFiles(path, thresholdTime, results)
			if err != nil {
				return nil, err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				log.Error("Error getting file info:", err)
				return nil, err
			}

			// Calculate the threshold time
			currentTime := time.Now()
			threshold := currentTime.Add(-time.Duration(thresholdTime) * time.Minute)

			if info.ModTime().Before(threshold) {
				// If the file's modification time is before the threshold time, it's considered old
				results = append(results, FileResult{FilePath: path, IsOld: true})
			} else {
				// If the file's modification time is after the threshold time, it's considered new
				results = append(results, FileResult{FilePath: path, IsOld: false})
			}
		}
	}
	return results, nil
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
