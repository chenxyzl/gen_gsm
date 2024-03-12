package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestGetterSetter(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testFile := filepath.Clean(filepath.Join(currentDir, "../../", "model/example.go"))
	testOutFile := filepath.Clean(filepath.Join(currentDir, "../../", "model/example.gen.go"))

	// 解析源代码
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, testFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse source code: %v", err)
	}

	genFile := &ast.File{
		Name:  f.Name,
		Decls: []ast.Decl{},
	}

	//
	addImport(genFile, []string{})

	// 找到需要生成 getter 和 setter 的 struct 类型
	ast.Inspect(f, func(n ast.Node) bool {
		if spec, ok := n.(*ast.TypeSpec); ok {
			structType, ok := spec.Type.(*ast.StructType)
			if ok {
				//格式检查
				fields := checkModelStruct(spec.Name, structType)
				// 生成 getter 和 setter 方法
				generateGettersAndSetters(genFile, spec.Name, structType, fields)
			}
		}
		return true
	})

	outputFile, err := os.Create(testOutFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// 打印修改后的源代码
	err = printer.Fprint(outputFile, fset, genFile)
	if err != nil {
		log.Fatalf("Failed to print modified source code: %v", err)
	}
}

func checkModelStruct(structNameIdent *ast.Ident, structType *ast.StructType) (out []*ast.Field) {
	contain := false
	for _, field := range structType.Fields.List {
		//获取类型名
		starExpr, ok := getImplType(field.Type)
		typeName := starExpr.Name
		//申明检查
		if len(field.Names) > 1 { //
			panic(fmt.Sprintf("类型:%v,每行只能声明1个filed", structNameIdent))
		} else if len(field.Names) == 0 { //匿名类型必须只能是DirtyModel
			if typeName != "DirtyModel" {
				panic(fmt.Sprintf("类型:%v, 字段:%v 必须申明名字", structNameIdent, typeName))
			}
			if ok {
				panic(fmt.Sprintf("类型:%v, 字段:%v 只能为值类型(去掉*)", structNameIdent, typeName))
			}
			contain = true
		} else { //具名类型如果不是基本类型,则必须是指针
			if isBasicType(typeName) {
				if ok {
					panic(fmt.Sprintf("类型:%v, 字段:%v 基本类型必须为值类型", structNameIdent, typeName))
				}
			} else {
				if !ok {
					panic(fmt.Sprintf("类型:%v, 字段:%v 非基本类型必须申明为指针", structNameIdent, typeName))
				}
			}
			out = append(out, field)

			fieldName := field.Names[0].Name
			if ast.IsExported(fieldName) {
				panic(fmt.Sprintf("类型:%v, 字段:%v 必须为非导出的(即小写)", structNameIdent, typeName))
			}
		}
	}
	if !contain {
		panic(fmt.Sprintf("类型:%v, 必须包含DirtyModel", structNameIdent))
	}
	return out
}

func generateGettersAndSetters(file *ast.File, structTypeExpr *ast.Ident, structType *ast.StructType, fields []*ast.Field) *ast.File {
	for idx, field := range fields {
		//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
		fieldName := field.Names[0].Name
		_, isStruct := getImplType(field.Type)

		file.Decls = append(file.Decls, &ast.FuncDecl{
			Name: ast.NewIdent("Get" + cases.Title(language.Und).String(fieldName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: field.Type,
						},
					},
				},
			},
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("s")},
						Type:  &ast.StarExpr{X: structTypeExpr},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("s"),
								Sel: ast.NewIdent(fieldName),
							},
						},
					},
				},
			},
		})

		var setterBody []ast.Stmt
		if isStruct {
			setterBody = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent(fieldName),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("v"),
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent("UpdateDirty"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(idx),
							},
						},
					},
				},
				&ast.IfStmt{
					If:   0,
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "v"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: ast.NewIdent("v.SetSelfDirtyIdx"),
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.INT,
											Value: strconv.Itoa(idx),
										},
										&ast.BasicLit{
											Kind:  token.FUNC,
											Value: "s.UpdateDirty",
										},
									},
								},
							},
						},
					},
					Else: nil,
				},
			}
		} else {
			setterBody = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent(fieldName),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("v"),
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent("UpdateDirty"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(idx),
							},
						},
					},
				},
			}
		}

		file.Decls = append(file.Decls, &ast.FuncDecl{
			Name: ast.NewIdent("Set" + cases.Title(language.Und).String(fieldName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("v")},
							Type:  field.Type,
						},
					},
				},
			},
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("s")},
						Type:  &ast.StarExpr{X: structTypeExpr},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: setterBody,
			},
		})
	}
	return file
}

func addImport(genFile *ast.File, imports []string) {
	// 添加导入语句
	if len(imports) > 0 {
		//importSpecs := make([]*ast.ImportSpec, 0, len(imports))
		importSpecs := make([]ast.Spec, 0, len(imports))
		for _, importPath := range imports {
			importSpecs = append(importSpecs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(importPath),
				},
			})
		}
		genFile.Decls = append(genFile.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		})
	}
}

func isBasicType(typeStr string) bool {
	basicTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"string", "bool",
		"byte", "rune",
	}
	for _, t := range basicTypes {
		if typeStr == t {
			return true
		}
	}
	return false
}

func isBasicType1(expr ast.Expr) bool {
	a, _ := getImplType(expr)
	return isBasicType(a.Name)
}

// getImplType 获取具体类型
// @return 具体类型
// @return 是否指针
func getImplType(expr ast.Expr) (*ast.Ident, bool) {
	var typeName *ast.Ident
	if selectExpr, selectOk := expr.(*ast.SelectorExpr); selectOk {
		return selectExpr.Sel, false
	} else if starExpr, starOk := expr.(*ast.StarExpr); starOk {
		if innerSelectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
			return innerSelectorExpr.Sel, true
		} else {
			return starExpr.X.(*ast.Ident), true
		}
	} else {
		typeName = expr.(*ast.Ident)
		return typeName, false
	}
}
