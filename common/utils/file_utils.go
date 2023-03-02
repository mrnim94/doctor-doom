package utils

import (
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"

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
					files = append(files, file)
				}
			}
		}(i)
	}

	// Scan files using filepath.WalkDir
	err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			// Send file to a worker
			worker := getWorkerNum(path, numWorkers)
			fileChans[worker] <- path
		}
		return nil
	})

	// Close channels to signal workers that no more files will be sent
	for i := range fileChans {
		close(fileChans[i])
	}

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
		return false
	}
	isAllFileName := nameMatch == "" || nameMatch == "*"
	if useAndOperator {
		return info.ModTime().Unix()*1000 >= ageMs && info.Size() >= sizeB && (isAllFileName || regexp.MustCompile(nameMatch).MatchString(info.Name()))
	} else {
		return info.ModTime().Unix()*1000 >= ageMs || info.Size() >= sizeB || (isAllFileName || regexp.MustCompile(nameMatch).MatchString(info.Name()))
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
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
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
