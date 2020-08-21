package main

import (
	"go/doc"
	"go/token"
)

// CopyPackage produces a json-annotated Package object from a GoDoc Package object.
func CopyPackage(pkg *doc.Package, fileSet *token.FileSet) (Package, error) {
	newPkg := Package{
		Type:       "package",
		Doc:        pkg.Doc,
		Name:       pkg.Name,
		ImportPath: pkg.ImportPath,
		Imports:    pkg.Imports,
		Filenames:  pkg.Filenames,
		Bugs:       pkg.Bugs,
	}

	newPkg.Notes = map[string][]*Note{}
	for key, value := range pkg.Notes {
		notes := make([]*Note, len(value))
		for i, note := range value {
			notes[i] = &Note{
				Pos:  note.Pos,
				End:  note.End,
				UID:  note.UID,
				Body: note.Body,
			}
		}
		newPkg.Notes[key] = notes
	}

	newPkg.Consts = CopyValues(pkg.Consts, pkg.Name, pkg.ImportPath, fileSet)
	var err error
	newPkg.Funcs, err = CopyFuncs(pkg.Funcs, pkg.Name, pkg.ImportPath, fileSet)
	if err != nil {
		return Package{}, err
	}

	newPkg.Types = make([]*Type, len(pkg.Types))
	for i, t := range pkg.Types {
		funcs, err := CopyFuncs(t.Funcs, pkg.Name, pkg.ImportPath, fileSet)
		if err != nil {
			return Package{}, err
		}

		methods, err := CopyFuncs(t.Methods, pkg.Name, pkg.ImportPath, fileSet)
		if err != nil {
			return Package{}, err
		}

		newPkg.Types[i] = &Type{
			Name:              t.Name,
			PackageName:       pkg.Name,
			PackageImportPath: pkg.ImportPath,
			Type:              "type",
			Consts:            CopyValues(t.Consts, pkg.Name, pkg.ImportPath, fileSet),
			Doc:               t.Doc,
			Funcs:             funcs,
			Methods:           methods,
			Vars:              CopyValues(t.Vars, pkg.Name, pkg.ImportPath, fileSet),
		}
	}

	newPkg.Vars = CopyValues(pkg.Vars, pkg.Name, pkg.ImportPath, fileSet)
	return newPkg, nil
}
