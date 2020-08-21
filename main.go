package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
)

func runMain(directory, filterRegexp string) error {
	fileSet := token.NewFileSet()
	filter, err := GetExcludeFilter(filterRegexp)
	if err != nil {
		return err
	}
	pkgs, firstError := parser.ParseDir(fileSet, directory, filter, parser.ParseComments|parser.AllErrors)
	if firstError != nil {
		return firstError
	}
	if len(pkgs) > 1 {
		return errors.New("Multiple packages found in directory")
	}

	for _, pkg := range pkgs {
		docPkg := doc.New(pkg, directory, 0)
		cleanedPkg := CopyPackage(docPkg, fileSet)
		pkgJSON, err := json.MarshalIndent(cleanedPkg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", pkgJSON)
	}

	return nil
}

func main() {
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
		log.Fatal("Fatal: Please specify a target_directory.")
	}

	err := runMain(directory, filterRegexp)
	if err != nil {
		log.Fatal(err)
	}
}
