package doom

import (
	"os"

	"github.com/mrnim94/doctor-doom/common/utils"
)

type DoomDestroyRules struct {
	Age  string `yaml:"age"`  // 1h, 1d, 1w, 1m, 1y. Default is 6 days
	Name string `yaml:"name"` // file name wildcard or regex (https://golang.org/pkg/path/filepath/#Match)
	Size string `yaml:"size"` // 1B, 1KB, 1MB, 1GB, 1TB. Default is 1B (no limit)
}

type DoomOptions struct {
	DoomPath   string           `yaml:"doom_path"`   // Root path to doom
	Circle     string           `yaml:"circle"`      // 1h, 1d, 1w, 1m, 1y, interval between each doom, or cron expression (https://godoc.org/github.com/robfig/cron). Default is every 7 days using cron expression (Sunday).
	DoomExport string           `yaml:"doom_export"` // Export log path folder
	Rule       DoomDestroyRules `yaml:"rule"`        // Rule to destroy files
	RuleAnd    bool             `yaml:"rule_and"`    // If true, all rules must be satisfied to destroy a file. If false, only one rule must be satisfied to destroy a file. Default is false.
}

func DefaultDoomOptions() DoomOptions {
	return DoomOptions{
		DoomPath:   "/tmp/doom",
		Circle:     "*/1 * * * *", // Every minute
		DoomExport: "/var/log",    // Export log path folder file will be named as doom-*.log
		Rule: DoomDestroyRules{
			Age:  "30d",  // 30 days
			Size: "100M", // 100MB
			Name: "*",    // All files
		},
		RuleAnd: false,
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

	defaultOptions.RuleAnd = defaultOptions.RuleAnd || overrideOptions.RuleAnd

	return defaultOptions
}

func DoomOptionsFromEnv() DoomOptions {
	doomPath := os.Getenv("DOOM_PATH")
	circle := os.Getenv("CIRCLE")
	doomExport := os.Getenv("DOOM_EXPORT")
	age := os.Getenv("RULE_AGE")
	size := os.Getenv("RULE_SIZE")
	name := os.Getenv("RULE_NAME")
	ruleAnd := os.Getenv("RULE_AND") == "true"

	return DoomOptions{
		DoomPath:   doomPath,
		Circle:     circle,
		DoomExport: doomExport,
		Rule: DoomDestroyRules{
			Age:  age,
			Size: size,
			Name: name,
		},
		RuleAnd: ruleAnd,
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
