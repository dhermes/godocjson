package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
)

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

	fileSet := token.NewFileSet()
	pkgs, firstError := parser.ParseDir(fileSet, directory, GetExcludeFilter(filterRegexp), parser.ParseComments|parser.AllErrors)
	if firstError != nil {
		panic(firstError)
	}
	if len(pkgs) > 1 {
		panic("Multiple packages found in directory!\n")
	}
	for _, pkg := range pkgs {
		docPkg := doc.New(pkg, directory, 0)
		cleanedPkg := CopyPackage(docPkg, fileSet)
		pkgJSON, err := json.MarshalIndent(cleanedPkg, "", "  ")
		if err != nil {
			log.Fatalf("Failed to encode JSON: %v", err)
		}
		fmt.Printf("%s\n", pkgJSON)
	}
}
