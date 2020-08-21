package main

import (
	"go/ast"
	"go/doc"
	"go/token"
)

func processFuncDecl(d *ast.FuncDecl, fun *Func) {
	fun.Params = make([]FuncParam, 0)
	for _, f := range d.Type.Params.List {
		t := typeOf(f.Type)
		for _, name := range f.Names {
			fun.Params = append(fun.Params, FuncParam{
				Type: t,
				Name: name.String(),
			})
		}
	}
	fun.Results = make([]FuncParam, 0)
	if d.Type.Results != nil {
		for _, f := range d.Type.Results.List {
			t := typeOf(f.Type)
			if len(f.Names) == 0 {
				// For case func foo() Type
				fun.Results = append(fun.Results, FuncParam{
					Type: t,
				})
			} else {
				// For case func foo() (name, name Type)
				for _, name := range f.Names {
					fun.Results = append(fun.Results, FuncParam{
						Type: t,
						Name: name.String(),
					})
				}
			}
		}
	}
}

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
