package main

import (
	"encoding/json"
	"errors"
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
	directory, filterRegexp, err := SetupCmd()
	if err != nil {
		log.Fatal(err)
	}

	err = runMain(directory, filterRegexp)
	if err != nil {
		log.Fatal(err)
	}
}
