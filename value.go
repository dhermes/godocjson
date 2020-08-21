package main

import (
	"go/doc"
	"go/token"
)

// CopyValues produces a json-annotated array of Value objects from an array of GoDoc Value objects.
func CopyValues(c []*doc.Value, packageName string, packageImportPath string, fileSet *token.FileSet) []*Value {
	newConsts := make([]*Value, len(c))
	for i, c := range c {
		position := fileSet.Position(c.Decl.TokPos)
		newConsts[i] = &Value{
			Doc:               c.Doc,
			Names:             c.Names,
			PackageName:       packageName,
			PackageImportPath: packageImportPath,
			Type:              c.Decl.Tok.String(),
			Filename:          position.Filename,
			Line:              position.Line,
		}
	}
	return newConsts
}
