package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Struct struct {
	Pkg    string
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	Tag  string
}

func parse(src, typ string) Struct {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}
	var resSt Struct
	var name string
	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		x, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if x.Obj != nil && x.Obj.Kind == ast.Typ {
			decl := x.Obj.Decl
			xx, ok := decl.(*ast.TypeSpec)
			if !ok {
				return true
			}
			name = xx.Name.String()
			if name != typ {
				return true
			}
			st, ok := xx.Type.(*ast.StructType)
			if !ok {
				return true
			}
			var fields []Field
			for _, field := range st.Fields.List {
				name := field.Names[0].String()
				ee, ok := field.Type.(*ast.Ident)
				if !ok {
					return true
				}
				fields = append(fields, Field{
					Name: name,
					Type: ee.String(),
				})
			}
			if len(fields) > 0 {
				resSt = Struct{
					Name:   name,
					Fields: fields,
				}
			}
			return true
		}
		return true
	})
	return resSt
}
