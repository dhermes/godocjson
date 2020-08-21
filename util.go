package main

import (
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"strings"
)

func typeOf(astValue interface{}) (result string, err error) {
	var subResult string
	switch typed := astValue.(type) {
	case *ast.Ident:
		result = typed.String()
	case *ast.ArrayType:
		subResult, err = typeOf(typed.Elt)
		result = "[]" + subResult
	case *ast.Field:
		subResult, err = typeOf(typed.Type)
		result = typed.Names[0].Name + " " + subResult
	case *ast.StructType:
		fields := make([]string, typed.Fields.NumFields())
		for i, f := range typed.Fields.List {
			subResult, err = typeOf(f.Type)
			if err != nil {
				return
			}
			fields[i] = subResult
		}
		result = fmt.Sprintf("struct{%s}", strings.Join(fields, ","))
	case *ast.InterfaceType:
		methods := make([]string, typed.Methods.NumFields())
		for i, m := range typed.Methods.List {
			subResult, err = typeOf(m.Type)
			if err != nil {
				return
			}
			methods[i] = subResult
		}
		result = fmt.Sprintf("interface{%s}", strings.Join(methods, ","))
	case *ast.SelectorExpr:
		subResult, err = typeOf(typed.X)
		result = subResult + "." + typed.Sel.Name
	case *ast.Ellipsis:
		subResult, err = typeOf(typed.Elt)
		result = "..." + subResult
	case *ast.StarExpr:
		subResult, err = typeOf(typed.X)
		result = "*" + subResult
	case *ast.FuncType:
		params := make([]string, typed.Params.NumFields())
		for i, p := range typed.Params.List {
			subResult, err = typeOf(p.Type)
			if err != nil {
				return
			}
			params[i] = subResult
		}
		var results []string
		if typed.Results != nil {
			results = make([]string, typed.Results.NumFields())
			for i, r := range typed.Results.List {
				subResult, err = typeOf(r.Type)
				if err != nil {
					return
				}
				results[i] = subResult
			}
		}
		result = fmt.Sprintf("func(%s)%s", strings.Join(params, ","), strings.Join(results, ","))
	case *ast.MapType:
		subResult, err = typeOf(typed.Key)
		if err != nil {
			return
		}
		var secondResult string
		secondResult, err = typeOf(typed.Value)
		result = fmt.Sprintf("map [%s]%s", subResult, secondResult)
	case *ast.ChanType:
		if typed.Dir == ast.SEND {
			subResult, err = typeOf(typed.Value)
			result = fmt.Sprintf("chan<- %s", subResult)
		} else if typed.Dir == ast.RECV {
			subResult, err = typeOf(typed.Value)
			result = fmt.Sprintf("<-chan %s", subResult)
		} else {
			subResult, err = typeOf(typed.Value)
			result = fmt.Sprintf("chan %s", subResult)
		}
	default:
		err = fmt.Errorf("Unknown type %+v", typed)
	}
	return
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
