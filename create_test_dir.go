package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"os"
// )

// // Create a directory with the given name in the current directory.
// //
// // The directory is created with the permissions 0755.
// //
// // It will have 10 levels of subdirectories.
// //
// // It each subdirectory will have 100 files with random names. and random content size between 10K and 100K.
// func CreateDir() {
// 	dirName := "tmp"
// 	os.Mkdir(dirName, 0755)
// 	for i := 0; i < 10; i++ {
// 		subDirName := fmt.Sprintf("%s/%d", dirName, i)
// 		os.Mkdir(subDirName, 0755)
// 		for j := 0; j < 100; j++ {
// 			fileName := fmt.Sprintf("%s/%d", subDirName, rand.Intn(100000))
// 			file, err := os.Create(fileName)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer file.Close()
// 			fileSize := rand.Intn(100000) + 10000
// 			for k := 0; k < fileSize; k++ {
// 				file.WriteString("a")
// 			}
// 		}
// 	}
// }

// func main() {
// 	CreateDir()
// }
