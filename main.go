package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"io"
	"linkcheck/models"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO add tests.
// TODO add CMD.
// TODO make this a linter for megalinter.
// TODO add workflow
func main() {
	start := time.Now()
	var app = cli.NewApp()
	appInfo(app)
	setupCLICommands(app)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Now()
	log.Info("Time elapsed: " + end.Sub(start).String())
}

func setupCLICommands(app *cli.App) {
	app.Commands = []cli.Command{{
		Name:    "check",
		Aliases: []string{"c"},
		Usage:   "Go through files and check the links",
		Action: func(ctx *cli.Context) {
			configPath := ctx.Args().Get(0)
			if configPath == "" {
				configPath = "configuration/linkcheck.json"
			}
			loadConfiguration(configPath)
			setUpLogger()
			readmeFiles := extractReadmeFiles()
			models.GetFilesProcessorInstance().Process(readmeFiles)
		},
		OnUsageError: nil,
		Flags:        nil,
	}}
}

func setUpLogger() {
	outputPath := viper.GetString("output_path")
	logFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).
			Fatal("Failed to open log file.")
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
}

func appInfo(cli *cli.App) {
	cli.Name = "GoMDLinkCheck"
	cli.Usage = "A linter in Golang to verify Markdown links."
	cli.Author = "Shiran Rubin"
	cli.Version = "1.0.0"
}

func loadConfiguration(configPath string) {
	viper.SetConfigName(configPath)
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	viper.AddConfigPath(basepath + "/configuration")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
	}
}

func extractReadmeFiles() []string {
	path := viper.GetString("path")
	var readmeFiles []string
	if envPath := os.Getenv("PROJECT_PATH"); envPath != "" {
		path = envPath
	}
	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() && strings.Contains(file.Name(), "vendor") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
			path, _ = filepath.Abs(path)
			readmeFiles = append(readmeFiles, path)
		}
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get files")
	}
	return readmeFiles
}
