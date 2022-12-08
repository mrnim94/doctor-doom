package main

import (
	doomcli "github.com/mrnim94/doctor-doom/doom_cli"
)

func main() {
	doomCli := doomcli.DoomCli{}
	doomCli.New()
	err := doomCli.Start()
	if err != nil {
		panic(err)
	}
}
