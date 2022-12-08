package utils

import (
	"os"
	"testing"
)

func prepare() {
	// Create folder structure like this:
	// ./tmp
	// - text.txt
	// - /folder1
	//   - text.txt
	//   - /folder2
	//     - text.txt

	// Create root folder
	os.Mkdir("./tmp", 0755)

	// Create file in root folder
	os.Create("./tmp/text.txt")

	// Create folder1
	os.Mkdir("./tmp/folder1", 0755)

	// Create file in folder1
	os.Create("./tmp/folder1/text.txt")

	// Create folder2
	os.Mkdir("./tmp/folder1/folder2", 0755)

	// Create file in folder2
	os.Create("./tmp/folder1/folder2/text.txt")
}

func cleanup() {
	// Remove folder structure
	os.RemoveAll("./tmp")
}

func TestListAllFiles(t *testing.T) {
	rootPath := "./tmp"
	recursive := true

	prepare()
	defer cleanup()

	fileUtils := FileUtils{}
	files := fileUtils.ListAllFiles(rootPath, recursive)
	if len(files) != 3 {
		t.Errorf("Expected 4 files, got %d", len(files))
	}
}

func BenchmarkListAllFiles(b *testing.B) {
	rootPath := "./tmp"
	recursive := true

	prepare()
	defer cleanup()

	fileUtils := FileUtils{}
	for i := 0; i < b.N; i++ {
		fileUtils.ListAllFiles(rootPath, recursive)
	}
}

func TestGetFileLastModifiedTime(t *testing.T) {
	filePath := "./tmp/text.txt"

	prepare()
	defer cleanup()

	fileUtils := FileUtils{}
	lastModifiedTime := fileUtils.GetFileLastModifiedTime(filePath)
	if lastModifiedTime == 0 {
		t.Errorf("Expected last modified time, got %d", lastModifiedTime)
	}
}

func BenchmarkGetFileLastModifiedTime(b *testing.B) {
	filePath := "./tmp/text.txt"

	prepare()
	defer cleanup()

	fileUtils := FileUtils{}
	for i := 0; i < b.N; i++ {
		fileUtils.GetFileLastModifiedTime(filePath)
	}
}

func TestParseYamlFile(t *testing.T) {
	filePath := "../../sample/config.yaml"

	fileUtils := FileUtils{}
	var value interface{}

	err := fileUtils.ParseYamlFile(filePath, &value)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func BenchmarkParseYamlFile(b *testing.B) {
	filePath := "../../sample/config.yaml"

	var value interface{}

	fileUtils := FileUtils{}

	for i := 0; i < b.N; i++ {
		fileUtils.ParseYamlFile(filePath, &value)
	}
}
