package doom

import (
	"errors"
	"fmt"
)

type DoomCommandBuilder struct {
	FolderPath     string
	NamePattern    string
	TargetSize     int64 // In bytes
	TargetLiveTime int64 // In ms
}

func (db *DoomCommandBuilder) WithNamePattern(namePattern string) *DoomCommandBuilder {
	db.NamePattern = namePattern
	return db
}

func (db *DoomCommandBuilder) WithTargetSize(size int64) *DoomCommandBuilder {
	db.TargetSize = size
	return db
}

func (db *DoomCommandBuilder) WithTargetLiveTime(timeMs int64) *DoomCommandBuilder {
	db.TargetLiveTime = timeMs
	return db
}

func (db *DoomCommandBuilder) WithFolderPath(folderPath string) *DoomCommandBuilder {
	db.FolderPath = folderPath
	return db
}

func (db *DoomCommandBuilder) bytesToMBs(bytes int64) int64 {
	return bytes / 1024 / 1024
}

func (db *DoomCommandBuilder) msToMins(ms int64) int64 {
	return ms / 1000 / 60
}

func (db *DoomCommandBuilder) BuildCommand() (string, error) {
	if db.FolderPath == "" {
		return "", errors.New("[DoctorDoomError] -- No target folder")
	}

	command := "find "

	command += fmt.Sprintf("%v -type f ", db.FolderPath)

	mb := db.bytesToMBs(db.TargetSize)
	command += fmt.Sprintf("-size +%vM ", mb)

	min := db.msToMins(db.TargetLiveTime)
	command += fmt.Sprintf("-mmin +%v ", min)

	// Delete command
	command += "-exec echo \"Deleting: {}\" \\; -exec rm {} \\;"

	return command, nil
}
