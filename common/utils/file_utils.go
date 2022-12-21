package utils

import (
	"os"
	"path/filepath"
	"regexp"
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
func (f *FileUtils) ListAllFilesMatch(rootPath string, ageMs int64, sizeB int64, nameMatch string, useAndOperator bool) []string {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if useAndOperator {
				fileLiveTimeMs := time.Now().Unix() - info.ModTime().Unix()
				fileName := info.Name()
				fileSize := info.Size()

				if fileLiveTimeMs >= ageMs && fileSize >= sizeB && f.MatchName(fileName, nameMatch) {
					files = append(files, path)
				}
			} else {
				if ageMs > 0 {
					nowUnix := time.Now().Unix()
					if nowUnix-info.ModTime().Unix() >= ageMs {
						files = append(files, path)
					}
				}
				if sizeB >= 0 {
					if info.Size() >= sizeB {
						files = append(files, path)
					}
				}
				if nameMatch != "" {
					if f.MatchName(info.Name(), nameMatch) {
						files = append(files, path)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
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

	return fileInfo.ModTime().Unix()
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
