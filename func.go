package main

import (
	"go/doc"
	"go/token"
)

// CopyFuncs produces a json-annotated array of Func objects from an array of GoDoc Func objects.
func CopyFuncs(f []*doc.Func, packageName string, packageImportPath string, fileSet *token.FileSet) []*Func {
	newFuncs := make([]*Func, len(f))
	for i, n := range f {
		position := fileSet.Position(n.Decl.Pos())
		newFuncs[i] = &Func{
			Doc:               n.Doc,
			Name:              n.Name,
			PackageName:       packageName,
			PackageImportPath: packageImportPath,
			Type:              "func",
			Orig:              n.Orig,
			Recv:              n.Recv,
			Filename:          position.Filename,
			Line:              position.Line,
		}
		processFuncDecl(n.Decl, newFuncs[i])
	}
	return newFuncs
}
