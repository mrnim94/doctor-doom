package doom

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mrnim94/doctor-doom/common/logger"
	"github.com/mrnim94/doctor-doom/common/utils"
	"github.com/robfig/cron/v3"
)

type DoctorDoom struct {
	DoomOptions DoomOptions
}

// Crate new DoctorDoom with the merge of default options, options from environment variables and options from arguments
func (doom *DoctorDoom) New(options DoomOptions) DoctorDoom {
	doomOptionsDefault := DefaultDoomOptions()
	doomOptionsFromEnv := DoomOptionsFromEnv()
	doomOptions := OverrideDoomOptions(doomOptionsDefault, doomOptionsFromEnv)
	doomOptions = OverrideDoomOptions(doomOptions, options)

	fmt.Println("Run with options: ", doomOptions)
	doom.DoomOptions = doomOptions
	return *doom
}

// Convert a list of files to a list of DoomVictim
func (doom *DoctorDoom) filesToDoomVictims(files []string) []DoomVictim {
	doomVictims := []DoomVictim{}
	fileUtils := utils.FileUtils{}

	for _, file := range files {
		doomVictims = append(doomVictims, DoomVictim{Path: file,
			Name:             strings.Split(file, "/")[len(strings.Split(file, "/"))-1],
			LastModifiedUnix: fileUtils.GetFileLastModifiedTime(file),
			Size:             fileUtils.GetFileSize(file),
		})
	}

	return doomVictims
}

// Seek for doom victims in the [doom_path] a.k.a [DOOM_PATH] environment variable
func (doom *DoctorDoom) GetDoomVictims() []DoomVictim {
	fileUtils := utils.FileUtils{}

	allFiles := fileUtils.ListAllFilesMatch(
		doom.DoomOptions.DoomPath,
		int64(doom.ageToMs(doom.DoomOptions.Rule.Age)),
		int64(doom.sizeToB(doom.DoomOptions.Rule.Size)),
		doom.DoomOptions.Rule.Name,
	)
	uniqueFiles := utils.ListToUnique(allFiles)
	doomVictims := doom.filesToDoomVictims(uniqueFiles)

	fmt.Println("Doom find", len(doomVictims), "doom victims", doomVictims)
	return doomVictims
}

// Destroy doom victims
func (doom *DoctorDoom) DestroyDoomVictims(doomVictims []DoomVictim) {
	fileUtils := utils.FileUtils{}
	doomLogger := logger.DoomLogger{}
	doomLogger.New(doom.DoomOptions.DoomExport)
	doomLogger.Info("Doom found", strconv.Itoa(len(doomVictims))+" doom victims"+fmt.Sprint(doomVictims))

	for _, doomVictim := range doomVictims {
		err := fileUtils.RemoveFile(doomVictim.Path)
		if err != nil {
			doomLogger.ErrorVictim("Doom destroy victim failed", doomVictim.Path, doomVictim.LastModifiedUnix, doomVictim.Size)
		}
		doomLogger.InfoVictim("Doom destroy", doomVictim.Path, doomVictim.LastModifiedUnix, doomVictim.Size)
	}
}

// Convert size to bytes
//
// Example:
//
//	1s -> 1000
//	1m -> 60000
//	1h -> 3600000
//	1d -> 86400000
//	1w -> 604800000
//	1M -> 2592000000
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

// Convert size to bytes
//
// Example:
//
//	1B -> 1
//	1K -> 1024
//	1M -> 1048576
//	1G -> 1073741824
//	1T -> 1099511627776
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

// Main function to destroy doom victims
func (doom *DoctorDoom) Destroy() {
	doomVictims := doom.GetDoomVictims()
	doom.DestroyDoomVictims(doomVictims)
}

var forever = make(chan bool)

func (doom *DoctorDoom) StartConquer() {
	fmt.Println("Start conquer the world 🌋")
	cron := cron.New()
	cron.AddFunc(doom.DoomOptions.Circle, func() {
		doom.Destroy()
	})
	cron.Start()

	// Block forever
	<-forever
}
