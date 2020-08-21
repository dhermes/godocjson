package main

import (
	"fmt"
	"go/ast"
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
