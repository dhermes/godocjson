package main

import (
	"go/token"
)

// Func represents a function declaration.
type Func struct {
	Doc               string      `json:"doc"`
	Name              string      `json:"name"`
	PackageName       string      `json:"packageName"`
	PackageImportPath string      `json:"packageImportPath"`
	Type              string      `json:"type"`
	Filename          string      `json:"filename"`
	Line              int         `json:"line"`
	Params            []FuncParam `json:"parameters"`
	Results           []FuncParam `json:"results"`

	// methods
	// (for functions, these fields have the respective zero value)
	Recv string `json:"recv"` // actual   receiver "T" or "*T"
	Orig string `json:"orig"` // original receiver "T" or "*T"
	// Level int    // embedding level; 0 means not embedded
}

// Package represents a package declaration.
type Package struct {
	Type       string             `json:"type"`
	Doc        string             `json:"doc"`
	Name       string             `json:"name"`
	ImportPath string             `json:"importPath"`
	Imports    []string           `json:"imports"`
	Filenames  []string           `json:"filenames"`
	Notes      map[string][]*Note `json:"notes"`
	// DEPRECATED. For backward compatibility Bugs is still populated,
	// but all new code should use Notes instead.
	Bugs []string `json:"bugs"`

	// declarations
	Consts []*Value `json:"consts"`
	Types  []*Type  `json:"types"`
	Vars   []*Value `json:"vars"`
	Funcs  []*Func  `json:"funcs"`
}

// Note represents a note comment.
type Note struct {
	Pos  token.Pos `json:"pos"`
	End  token.Pos `json:"end"`  // position range of the comment containing the marker
	UID  string    `json:"uid"`  // uid found with the marker
	Body string    `json:"body"` // note body text
}

// Type represents a type declaration.
type Type struct {
	PackageName       string `json:"packageName"`
	PackageImportPath string `json:"packageImportPath"`
	Doc               string `json:"doc"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Filename          string `json:"filename"`
	Line              int    `json:"line"`
	// Decl              *ast.GenDecl

	// associated declarations
	Consts  []*Value `json:"consts"`  // sorted list of constants of (mostly) this type
	Vars    []*Value `json:"vars"`    // sorted list of variables of (mostly) this type
	Funcs   []*Func  `json:"funcs"`   // sorted list of functions returning this type
	Methods []*Func  `json:"methods"` // sorted list of methods (including embedded ones) of this type
}

// Value represents a value declaration.
type Value struct {
	PackageName       string   `json:"packageName"`
	PackageImportPath string   `json:"packageImportPath"`
	Doc               string   `json:"doc"`
	Names             []string `json:"names"` // var or const names in declaration order
	Type              string   `json:"type"`
	Filename          string   `json:"filename"`
	Line              int      `json:"line"`
	// Decl              *ast.GenDecl
}

// FuncParam represents a parameter to a function.
type FuncParam struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
