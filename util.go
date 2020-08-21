package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"os"
	"regexp"
	"strings"
)

func typeOf(x interface{}) string {
	switch x := x.(type) {
	case *ast.Ident:
		return x.String()
	case *ast.ArrayType:
		return "[]" + typeOf(x.Elt)
	case *ast.Field:
		return x.Names[0].Name + " " + typeOf(x.Type)
	case *ast.StructType:
		fields := make([]string, x.Fields.NumFields())
		for i, f := range x.Fields.List {
			fields[i] = typeOf(f.Type)
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, ","))
	case *ast.InterfaceType:
		methods := make([]string, x.Methods.NumFields())
		for i, m := range x.Methods.List {
			methods[i] = typeOf(m.Type)
		}
		return fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
	case *ast.SelectorExpr:
		return typeOf(x.X) + "." + x.Sel.Name
	case *ast.Ellipsis:
		return "..." + typeOf(x.Elt)
	case *ast.StarExpr:
		return "*" + typeOf(x.X)
	case *ast.FuncType:
		params := make([]string, x.Params.NumFields())
		for i, p := range x.Params.List {
			params[i] = typeOf(p.Type)
		}
		var results []string
		if x.Results != nil {
			results = make([]string, x.Results.NumFields())
			for i, r := range x.Results.List {
				results[i] = typeOf(r.Type)
			}
		}
		return fmt.Sprintf("func(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
	case *ast.MapType:
		return fmt.Sprintf("map [%s]%s", typeOf(x.Key), typeOf(x.Value))
	case *ast.ChanType:
		if x.Dir == ast.SEND {
			return fmt.Sprintf("chan<- %s", typeOf(x.Value))
		} else if x.Dir == ast.RECV {
			return fmt.Sprintf("<-chan %s", typeOf(x.Value))
		} else {
			return fmt.Sprintf("chan %s", typeOf(x.Value))
		}
	default:
		panic(fmt.Sprintf("Unknown type %+v", x))
	}
}

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

// CopyPackage produces a json-annotated Package object from a GoDoc Package object.
func CopyPackage(pkg *doc.Package, fileSet *token.FileSet) Package {
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
	newPkg.Funcs = CopyFuncs(pkg.Funcs, pkg.Name, pkg.ImportPath, fileSet)

	newPkg.Types = make([]*Type, len(pkg.Types))
	for i, t := range pkg.Types {
		newPkg.Types[i] = &Type{
			Name:              t.Name,
			PackageName:       pkg.Name,
			PackageImportPath: pkg.ImportPath,
			Type:              "type",
			Consts:            CopyValues(t.Consts, pkg.Name, pkg.ImportPath, fileSet),
			Doc:               t.Doc,
			Funcs:             CopyFuncs(t.Funcs, pkg.Name, pkg.ImportPath, fileSet),
			Methods:           CopyFuncs(t.Methods, pkg.Name, pkg.ImportPath, fileSet),
			Vars:              CopyValues(t.Vars, pkg.Name, pkg.ImportPath, fileSet),
		}
	}

	newPkg.Vars = CopyValues(pkg.Vars, pkg.Name, pkg.ImportPath, fileSet)
	return newPkg
}

// GetExcludeFilter builds a filter function that can be used with
// `parser.ParseDir`.
func GetExcludeFilter(re string) (func(os.FileInfo) bool, error) {
	if re == "" {
		return nil, nil
	}

	pattern, err := regexp.Compile(re)
	if err != nil {
		return nil, err
	}

	filter := func(info os.FileInfo) bool {
		matched := pattern.MatchString(info.Name())
		return !matched
	}
	return filter, nil
}
