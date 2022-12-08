package doom

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mrnim94/doctor-doom/common/utils"
)

type DoctorDoom struct {
	DoomOptions DoomOptions
}

func (doom *DoctorDoom) New(options DoomOptions) DoctorDoom {
	doomOptionsDefault := DefaultDoomOptions()
	doomOptionsFromEnv := DoomOptionsFromEnv()
	doomOptions := OverrideDoomOptions(doomOptionsDefault, doomOptionsFromEnv)
	doomOptions = OverrideDoomOptions(doomOptions, options)

	doom.DoomOptions = doomOptions
	return *doom
}

func (doom *DoctorDoom) filesToDoomVictims(files []string) []DoomVictim {
	doomVictims := []DoomVictim{}
	fileUtils := utils.FileUtils{}
	for _, file := range files {
		doomVictims = append(doomVictims, DoomVictim{Path: file,
			Name:             strings.Split(file, "/")[len(strings.Split(file, "/"))-1],
			LastModifiedUnix: fileUtils.GetFileLastModifiedTime(file),
		})
	}
	return doomVictims
}

func (doom *DoctorDoom) GetDoomVictims() []DoomVictim {
	fileUtils := utils.FileUtils{}
	allFiles := fileUtils.ListAllFilesMatch(doom.DoomOptions.DoomPath,
		int64(doom.ageToMs(doom.DoomOptions.Rule.Age)),
		int64(doom.sizeToB(doom.DoomOptions.Rule.Size)), doom.DoomOptions.Rule.Name)
	uniqueFiles := utils.ListToUnique(allFiles)
	doomVictims := doom.filesToDoomVictims(uniqueFiles)
	return doomVictims
}

func (doom *DoctorDoom) DestroyDoomVictims(doomVictims []DoomVictim) {
	fileUtils := utils.FileUtils{}
	for _, doomVictim := range doomVictims {
		err := fileUtils.RemoveFile(doomVictim.Path)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (doom *DoctorDoom) ageToMs(age string) int {
	unit := age[len(age)-1:]
	ageInt, err := strconv.Atoi(age[:len(age)-1])
	if err != nil {
		fmt.Println(err)
		return 0
	}

	switch unit {
	case "s":
		return ageInt * 1000
	case "m":
		return ageInt * 1000 * 60
	case "h":
		return ageInt * 1000 * 60 * 60
	case "d":
		return ageInt * 1000 * 60 * 60 * 24
	case "w":
		return ageInt * 1000 * 60 * 60 * 24 * 7
	case "M":
		return ageInt * 1000 * 60 * 60 * 24 * 30
	case "y":
		return ageInt * 1000 * 60 * 60 * 24 * 365
	default:
		return 0
	}
}

func (doom *DoctorDoom) sizeToB(size string) int {
	unit := size[len(size)-1:]
	sizeInt, err := strconv.Atoi(size[:len(size)-1])
	if err != nil {
		fmt.Println(err)
		return 0
	}

	switch unit {
	case "B":
		return sizeInt
	case "K":
		return sizeInt * 1024
	case "M":
		return sizeInt * 1024 * 1024
	case "G":
		return sizeInt * 1024 * 1024 * 1024
	case "T":
		return sizeInt * 1024 * 1024 * 1024 * 1024
	default:
		return 0
	}
}

func (doom *DoctorDoom) Destroy() {
	doomVictims := doom.GetDoomVictims()
	doom.DestroyDoomVictims(doomVictims)
}
