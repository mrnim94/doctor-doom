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

		//// Execute the find command
		//var cmd *exec.Cmd
		//// Check the operating system and execute the appropriate command
		//if runtime.GOOS == "windows" {
		//	cmd = exec.Command("cmd", "/c", "dir", "/b", "/s", rootPath)
		//} else {
		//	cmd = exec.Command("find", rootPath, "-type", "f", "-mmin", "+"+strconv.FormatInt(minutes, 10))
		//}
		//log.Debug("Running command: ", cmd)
		//output, err := cmd.StdoutPipe()
		//if err != nil {
		//	log.Error("Error creating stdout pipe:", err)
		//}
		//
		//stderr, err := cmd.StderrPipe()
		//if err != nil {
		//	log.Error("Error creating stderr pipe:", err)
		//}
		//
		//if err := cmd.Start(); err != nil {
		//	log.Error("Error starting command:", err)
		//}
		//
		//// Debug: Print error output
		//scannerErr := bufio.NewScanner(stderr)
		//for scannerErr.Scan() {
		//	log.Error("STDERR:", scannerErr.Text())
		//}

		files, err := os.ReadDir(rootPath)
		if err != nil {
			log.Error(err)
		}

		results := make(chan FileResult)
		var wg sync.WaitGroup

		log.Debug("Begin to Check Old File")
		//scanner := bufio.NewScanner(output)
		fileCount := 0
		for _, file := range files {

			wg.Add(1)

			go func(filePath string) {
				defer wg.Done()
				isOld, err := checkOlFile(filePath, minutes)
				if err != nil {
					log.Error("Error checking file:", err)
					return
				}
				results <- FileResult{FilePath: filePath, IsOld: isOld}
			}(filepath.Join(rootPath, file.Name()))
			fileCount++
		}

		log.Debug("The number of files that are found is: ", fileCount)

		go func() {
			wg.Wait()
			close(results)
		}()

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
