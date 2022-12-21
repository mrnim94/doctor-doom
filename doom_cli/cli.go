package doomcli

import (
	"fmt"
	"os"

	"github.com/mrnim94/doctor-doom/common/utils"
	"github.com/mrnim94/doctor-doom/doom"
	"github.com/urfave/cli/v2"
)

type DoomCli struct {
	cliApp *cli.App
}

var doomCliApp *cli.App

func init() {
	doomCliApp = &cli.App{
		Name:      "doom",
		Usage:     "Doctor Doom. Conquer the world, destroy victims 🔥🔥🔥",
		UsageText: "doom [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "doom-path",
				Usage: "The path where Doctor Doom will seek for doom victims. MUST BE ABSOLUTE PATH.",
			},
			&cli.StringFlag{
				Name:  "doom-export",
				Usage: "The path where Doctor Doom will export the log file.",
			},
			&cli.StringFlag{
				Name:  "age",
				Usage: "The age of the doom victims. Format: 1h, 1d, 1w, 1m, 1y",
			},
			&cli.StringFlag{
				Name:  "size",
				Usage: "The size of the doom victims. Format: 1B, 1K, 1M, 1G, 1T",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "The name of the doom victims. Format: (/s or *) for all files, regular expression for specific files",
			},
			&cli.StringFlag{
				Name:  "doom-config",
				Usage: "The path to the doom config file. Values in the config file will be overwritten by the cli flags",
			},
			&cli.StringFlag{
				Name:  "doom-circle",
				Usage: "The circle of the doom. Cron expression (https://godoc.org/github.com/robfig/cron)",
			},
			&cli.BoolFlag{
				Name:  "rule-and",
				Usage: "The rule of the doom. If true, all rules must be satisfied. If false, at least one rule must be satisfied",
			},
		},
		Action: doomCliAction,
	}
}

func doomCliAction(c *cli.Context) error {
	fmt.Println("Conquer the world, destroy victims 🔥🔥🔥")
	var doomOptions doom.DoomOptions

	// Get doom config
	doomConfig := c.String("doom-config")
	if doomConfig != "" {
		fmt.Println("Doom config: ", doomConfig)
		var fileUtils utils.FileUtils
		err := fileUtils.ParseYamlFile(doomConfig, &doomOptions)
		if err != nil {
			fmt.Println("Error parsing doom config file: ", err)
			doomOptions = doom.DefaultDoomOptions()
		}
	} else {
		// Get doom path
		doomPath := c.String("doom-path")
		// if doomPath == "" {
		// 	panic("Doom path is required")
		// }
		// Get doom export
		doomExport := c.String("doom-export")

		// Get age
		age := c.String("age")

		// Get size
		size := c.String("size")

		// Get name
		name := c.String("name")

		// Get circle
		circle := c.String("doom-circle")

		// Rule and
		ruleAnd := c.Bool("rule-and")

		doomOptions = doom.DoomOptions{
			DoomPath:   doomPath,
			DoomExport: doomExport,
			Rule: doom.DoomDestroyRules{
				Age:  age,
				Size: size,
				Name: name,
			},
			Circle:  circle,
			RuleAnd: ruleAnd,
		}
	}

	// Start conquer the world
	doctorDoom := doom.DoctorDoom{}
	doctorDoom.New(doomOptions)
	doctorDoom.StartConquer()
	return nil
}

// Create new DoomCli
func (doomCli *DoomCli) New() DoomCli {
	doomCli.cliApp = doomCliApp
	return *doomCli
}

func (doomCli *DoomCli) Start() error {
	return doomCli.cliApp.Run(os.Args)
}
