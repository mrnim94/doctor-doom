package handler

import (
	"doctor_doom/helper"
	"doctor_doom/log"
	"github.com/go-co-op/gocron"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	s := gocron.NewScheduler(time.UTC)

	deleteTask := func() {
		rootPath := helper.GetEnvOrDefault("DOOM_PATH", "test_delete")

		minutes := helper.DurationToMinutes(helper.GetEnvOrDefault("RULE_AGE", "1m"))

		results := make(chan FileResult)
		var wg sync.WaitGroup

		log.Debug("Begin to Check Old File")

		// Start the recursive file listing and processing
		wg.Add(1)
		go listFiles(rootPath, minutes, &wg, results)

		// Close the results channel once all processing is done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Handle the results
		for result := range results {
			if result.IsOld {
				log.Debug("File ", result.FilePath, " is older than threshold")
				go deleteFile(result.FilePath)
			} else {
				log.Debug("File ", result.FilePath, " is not older than threshold")
			}
		}
	}

	_, err := s.Cron(helper.GetEnvOrDefault("CIRCLE", "*/1 * * * *")).Do(deleteTask)
	if err != nil {
		log.Error(err)
	}
	s.StartAsync()
}

func listFiles(dir string, thresholdTime int64, wg *sync.WaitGroup, results chan<- FileResult) {
	defer wg.Done()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			// It's a directory, recurse into it
			wg.Add(1)
			go listFiles(path, thresholdTime, wg, results)
		} else {
			// It's a file, process it concurrently
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				isOld, err := checkOlFile(filePath, thresholdTime)
				if err != nil {
					log.Error("Error checking file:", err)
					return
				}
				results <- FileResult{FilePath: filePath, IsOld: isOld}
			}(path)
		}
	}
}

func checkOlFile(filePath string, thresholdTime int64) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Error(err)
	}
	modTime := fileInfo.ModTime()
	currentTime := time.Now()
	thresholdDuration := time.Duration(thresholdTime) * time.Minute
	return modTime.Add(thresholdDuration).Before(currentTime), nil

}

func deleteFile(filePath string) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// On Windows, use the "cmd" shell to capture output
		cmd = exec.Command("cmd", "/c", "del", filePath)
	} else {
		// On Linux and other Unix-based systems, use "sh" to capture output
		cmd = exec.Command("sh", "-c", "rm "+filePath)
	}

	// Capture and print the output
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Error("Error deleting the file:", err)
	}
	log.Debug("Running command: ", cmd)
	log.Debug("Command Output: ", strings.TrimSpace(string(output)))
}
