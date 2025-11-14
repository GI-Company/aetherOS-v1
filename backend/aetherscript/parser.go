
package aetherscript

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Parse takes a string of AetherScript code and returns an AST.
func Parse(code string) (*ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AetherScript: %w", err)
	}
	return file, nil
}
