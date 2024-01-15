package utils

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type FileUtils struct{}

// Return all files in a directory
//
// # If recursive is true, it will return all files in all level of subdirectories
//
// @param rootPath string
//
// @param recursive bool
func (f *FileUtils) ListAllFiles(rootPath string, recursive bool) []string {
	var files []string

	if recursive {
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			panic(err)
		}

	} else {
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
				return filepath.SkipDir
			}
			return nil
		})

		if err != nil {
			panic(err)
		}
	}
	return files
}

// Return all the files that match one or more conditions:
//
// - File last modified time is in Unix timestamp to now in Ms (e.g. 1 day = 86400000) is greater than ageMs
//
// - File size is greater than sizeB
//
// - File name matches the regex nameMatch
//
// @param rootPath string
//
// @param ageMs int64
//
// @param sizeB int64
//
// @param nameMatch string
func (f *FileUtils) ListAllFilesMatch(rootPath string, ageMs int64, sizeB int64, nameMatch string, useAndOperator bool, numWorkers int) []string {
	var files []string
	var mu sync.Mutex

	// Create a channel for each worker
	fileChans := make([]chan string, numWorkers)
	for i := range fileChans {
		fileChans[i] = make(chan string, 100)
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for file := range fileChans[i] {
				if shouldProcessFile(file, ageMs, sizeB, nameMatch, useAndOperator) {
					mu.Lock()
					files = append(files, file)
					mu.Unlock()
				}
			}
		}(i)
	}

	//// Scan files using filepath.WalkDir
	//err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	if !entry.IsDir() {
	//		// Send file to a worker
	//		worker := getWorkerNum(path, numWorkers)
	//		fileChans[worker] <- path
	//	}
	//	return nil
	//})
	//

	// Convert milliseconds to minutes
	minutes := int(ageMs / 1000 / 60)

	// Convert minutes to string and prepend with "+"
	minutesStr := "+" + strconv.Itoa(minutes)

	// Execute the find command
	var cmd *exec.Cmd
	// Check the operating system and execute the appropriate command
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "dir", "/b", "/s", rootPath)
	} else {
		cmd = exec.Command("find", rootPath, "-type", "f", "-mmin", minutesStr)
	}

	fmt.Printf("Running command: %v\n", cmd)

	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)

	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
	}

	// Debug: Print error output
	scannerErr := bufio.NewScanner(stderr)
	for scannerErr.Scan() {
		fmt.Println("STDERR:", scannerErr.Text())
	}

	fileCount := 0
	// Read and process the output
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		// Process each file path (in this case, just print it)
		worker := getWorkerNum(line, numWorkers)
		fileChans[worker] <- line
		fileCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
	}

	// Close channels to signal workers that no more files will be sent
	for i := range fileChans {
		close(fileChans[i])
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command to finish:", err)
	}

	// Print the number of files
	fmt.Println("Number of files:", fileCount)

	// Wait for workers to finish
	wg.Wait()

	if err != nil {
		return []string{}
	}

	return files

}

func shouldProcessFile(path string, ageMs int64, sizeB int64, nameMatch string, useAndOperator bool) bool {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("ERROR - shouldProcessFile ", err)
		return false
	}
	isAllFileName := nameMatch == "" || nameMatch == "*"
	now := time.Now().Unix() * 1000
	modTimeDiff := now - info.ModTime().Unix()*1000
	if useAndOperator {
		return modTimeDiff >= ageMs && info.Size() >= sizeB && (isAllFileName || regexp.MustCompile(nameMatch).MatchString(info.Name()))
	} else {
		return modTimeDiff >= ageMs || info.Size() >= sizeB || (isAllFileName || regexp.MustCompile(nameMatch).MatchString(info.Name()))
	}
}

func getWorkerNum(path string, numWorkers int) int {
	h := fnv.New32a()
	h.Write([]byte(path))
	return int(h.Sum32()) % numWorkers
}

// Do the regex match
func (f *FileUtils) MatchName(name string, match string) bool {
	if match == "" {
		return true
	}

	if match == "*" {
		return true
	}

	return regexp.MustCompile(match).MatchString(name)
}

// Return file last modified time in Unix timestamp
//
// @param filePath string
func (f *FileUtils) GetFileLastModifiedTime(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	return fileInfo.ModTime().Unix() * 1000
}

// Parse a YAML file into a struct
//
// @param filePath string
//
// @param target interface{}
func (f *FileUtils) ParseYamlFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

// Remove file
func (f *FileUtils) RemoveFile(filePath string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Using PowerShell to remove the file
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Remove-Item '%s'", filePath))
	} else {
		// Using 'rm' command for Unix-like systems
		cmd = exec.Command("rm", filePath)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	successMsg := fmt.Sprintf("Successfully deleted file using %s: %s\n", cmd.Path, filePath)
	fmt.Printf(successMsg)

	return nil
}

// Return file size in bytes
func (f *FileUtils) GetFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	return fileInfo.Size()
}
