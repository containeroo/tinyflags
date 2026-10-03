package tinyflags_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// TestProductionDeclarationsAreDocumented keeps the production code's local
// documentation convention intact: every function, struct, and struct field
// has a descriptive comment. Test files are intentionally excluded.
func TestProductionDeclarationsAreDocumented(t *testing.T) {
	root := filepath.Clean("..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == ".projectatlas" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		return verifyFileDocumentation(path)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// verifyFileDocumentation reports the first undocumented production declaration in path.
func verifyFileDocumentation(path string) error {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, declaration := range file.Decls {
		switch declaration := declaration.(type) {
		case *ast.FuncDecl:
			if declaration.Doc == nil {
				return undocumentedDeclaration(path, fileSet, declaration.Pos(), "function")
			}
		case *ast.GenDecl:
			for _, specification := range declaration.Specs {
				typeSpec, ok := specification.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				if declaration.Doc == nil && typeSpec.Doc == nil {
					return undocumentedDeclaration(path, fileSet, typeSpec.Pos(), "struct")
				}
				for _, field := range structType.Fields.List {
					if field.Doc == nil && field.Comment == nil {
						return undocumentedDeclaration(path, fileSet, field.Pos(), "struct field")
					}
				}
			}
		}
	}
	return nil
}

// undocumentedDeclaration formats an actionable documentation failure.
func undocumentedDeclaration(path string, fileSet *token.FileSet, position token.Pos, kind string) error {
	return fmt.Errorf("%s at %s requires a documentation comment", kind, fileSet.Position(position))
}
