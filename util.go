package main

import (
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"strings"
)

func typeOf(astValue interface{}) string {
	switch typed := astValue.(type) {
	case *ast.Ident:
		return typed.String()
	case *ast.ArrayType:
		return "[]" + typeOf(typed.Elt)
	case *ast.Field:
		return typed.Names[0].Name + " " + typeOf(typed.Type)
	case *ast.StructType:
		fields := make([]string, typed.Fields.NumFields())
		for i, f := range typed.Fields.List {
			fields[i] = typeOf(f.Type)
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, ","))
	case *ast.InterfaceType:
		methods := make([]string, typed.Methods.NumFields())
		for i, m := range typed.Methods.List {
			methods[i] = typeOf(m.Type)
		}
		return fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
	case *ast.SelectorExpr:
		return typeOf(typed.X) + "." + typed.Sel.Name
	case *ast.Ellipsis:
		return "..." + typeOf(typed.Elt)
	case *ast.StarExpr:
		return "*" + typeOf(typed.X)
	case *ast.FuncType:
		params := make([]string, typed.Params.NumFields())
		for i, p := range typed.Params.List {
			params[i] = typeOf(p.Type)
		}
		var results []string
		if typed.Results != nil {
			results = make([]string, typed.Results.NumFields())
			for i, r := range typed.Results.List {
				results[i] = typeOf(r.Type)
			}
		}
		return fmt.Sprintf("func(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
	case *ast.MapType:
		return fmt.Sprintf("map [%s]%s", typeOf(typed.Key), typeOf(typed.Value))
	case *ast.ChanType:
		if typed.Dir == ast.SEND {
			return fmt.Sprintf("chan<- %s", typeOf(typed.Value))
		} else if typed.Dir == ast.RECV {
			return fmt.Sprintf("<-chan %s", typeOf(typed.Value))
		} else {
			return fmt.Sprintf("chan %s", typeOf(typed.Value))
		}
	default:
		panic(fmt.Sprintf("Unknown type %+v", typed))
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
