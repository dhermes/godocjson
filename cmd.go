package main

import (
	"errors"
	"flag"
	"log"
)

// GetUsageText returns the usage text for the command.
func GetUsageText() {
	log.Println("Usage of godocjson:")
	log.Println("godocjson [-e] target_directory")
	flag.PrintDefaults()
}

// SetupCmd parses flags and does other stuff you would do with cobra.
func SetupCmd() (string, string, error) {
	var filterRegexp string
	// Disable timestamps inside the log file as we will just use it as wrapper
	// around stderr for now.
	log.SetFlags(0)

	flag.Usage = GetUsageText
	flag.StringVar(&filterRegexp, "e", "", "Regex filter for excluding source files")
	flag.Parse()

	directory := flag.Arg(0)
	if directory == "" {
		flag.Usage()
		return "", "", errors.New("Please specify a target_directory")
	}

	return directory, filterRegexp, nil
}
