package doom

import (
	"os"

	"github.com/mrnim94/doctor-doom/common/utils"
)

type DoomDestroyRules struct {
	Age  string `json:"age"`  // 1h, 1d, 1w, 1m, 1y. Default is 6 days
	Name string `json:"name"` // file name wildcard or regex (https://golang.org/pkg/path/filepath/#Match)
	Size string `json:"size"` // 1B, 1KB, 1MB, 1GB, 1TB. Default is 1B (no limit)
}

type DoomOptions struct {
	DoomPath   string           `json:"doom_path"`   // Root path to doom
	Circle     string           `json:"circle"`      // 1h, 1d, 1w, 1m, 1y, interval between each doom, or cron expression (https://godoc.org/github.com/robfig/cron). Default is every 7 days using cron expression (Sunday).
	DoomExport string           `json:"doom_export"` // Export log path folder
	Rule       DoomDestroyRules `json:"rule"`        // Rule to destroy files
}

func DefaultDoomOptions() DoomOptions {
	return DoomOptions{
		DoomPath:   "/tmp/doom",
		Circle:     "0 0 0 * * 0", // Every Sunday at 00:00:00
		DoomExport: "/var/log",    // Export log path folder file will be named as doom-*.log
		Rule: DoomDestroyRules{
			Age:  "30d",  // 6 days
			Size: "100M", // 1 byte
			Name: "*",    // All files
		},
	}
}

func OverrideDoomOptions(defaultOptions DoomOptions, overrideOptions DoomOptions) DoomOptions {
	if overrideOptions.DoomPath != "" {
		defaultOptions.DoomPath = overrideOptions.DoomPath
	}

	if overrideOptions.Circle != "" {
		defaultOptions.Circle = overrideOptions.Circle
	}

	if overrideOptions.DoomExport != "" {
		defaultOptions.DoomExport = overrideOptions.DoomExport
	}

	if overrideOptions.Rule.Age != "" {
		defaultOptions.Rule.Age = overrideOptions.Rule.Age
	}

	if overrideOptions.Rule.Size != "" {
		defaultOptions.Rule.Size = overrideOptions.Rule.Size
	}

	if overrideOptions.Rule.Name != "" {
		defaultOptions.Rule.Name = overrideOptions.Rule.Name
	}

	return defaultOptions
}

func DoomOptionsFromEnv() DoomOptions {
	doomPath := os.Getenv("DOOM_PATH")
	circle := os.Getenv("DOOM_CIRCLE")
	doomExport := os.Getenv("DOOM_EXPORT")
	age := os.Getenv("RULE_AGE")
	size := os.Getenv("RULE_SIZE")
	name := os.Getenv("RULE_NAME")

	return DoomOptions{
		DoomPath:   doomPath,
		Circle:     circle,
		DoomExport: doomExport,
		Rule: DoomDestroyRules{
			Age:  age,
			Size: size,
			Name: name,
		},
	}
}

func DoomOptionsFromConfigFile(configFile string) DoomOptions {
	fileUtils := utils.FileUtils{}
	var doomOptions DoomOptions

	err := fileUtils.ParseYamlFile(configFile, &doomOptions)
	if err != nil {
		panic(err)
	}

	return doomOptions
}
