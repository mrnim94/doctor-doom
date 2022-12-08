package main

import (
	"github.com/mrnim94/doctor-doom/common/logger"
	"github.com/mrnim94/doctor-doom/doom"
)

func main() {
	doctorDoom := doom.DoctorDoom{}
	doctorDoom.New(doom.DoomOptions{})

	logger.DoomLoggerInit(doctorDoom.DoomOptions.DoomExport)

	// Start conquer the world
	doctorDoom.StartConquer()
}
