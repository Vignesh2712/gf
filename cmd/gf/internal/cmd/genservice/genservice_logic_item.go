// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/gogf/gf/v2/os/gfile"
)

type logicItem struct {
	Receiver   string              `eg:"sUser"`
	MethodName string              `eg:"GetList"`
	Params     []map[string]string `eg:"ctx: context.Context, cond: *SearchInput"`
	Results    []map[string]string `eg:"list: []*User, err: error"`
	Comment    string              `eg:"Get user list"`
}

// CalculateItemsInSrc retrieves the logic items in the specified source file.
// It can't skip the private methods.
// It can't skip the imported packages of import alias equal to `_`.
func (c CGenService) CalculateItemsInSrc(filePath string) (pkgItems []packageItem, logicItems []logicItem, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			// calculate the imported packages
			pkgItems = append(pkgItems, c.getImportPackages(x))

		case *ast.FuncDecl:
			// calculate the logic items
			if x.Recv == nil {
				return true
			}

			var funcName = x.Name.Name
			logicItems = append(logicItems, logicItem{
				Receiver:   c.getFuncReceiverTypeName(x),
				MethodName: funcName,
				Params:     c.getFuncParams(x),
				Results:    c.getFuncResults(x),
				Comment:    c.getFuncComment(x),
			})
		}
		return true
	})
	return
}

// getImportPackages retrieves the imported packages from the specified ast.ImportSpec.
func (c CGenService) getImportPackages(node *ast.ImportSpec) (packages packageItem) {
	if node.Path == nil {
		return
	}
	var (
		alias     string
		path      = node.Path.Value
		rawImport string
	)
	if node.Name != nil {
		alias = node.Name.Name
		rawImport = alias + " " + path
	} else {
		rawImport = path
	}
	return packageItem{
		Alias:     alias,
		Path:      path,
		RawImport: rawImport,
	}
}

// getFuncReceiverTypeName retrieves the receiver type of the function.
// For example:
//
// func(s *sArticle) -> *sArticle
// func(s sArticle) -> sArticle
func (c CGenService) getFuncReceiverTypeName(node *ast.FuncDecl) (receiverType string) {
	if node.Recv == nil {
		return ""
	}
	receiverType, err := c.astExprToString(node.Recv.List[0].Type)
	if err != nil {
		return ""
	}
	return
}

// getFuncParams retrieves the input parameters of the function.
// It returns the name and type of the input parameters.
// For example:
//
// []map[string]string{paramName:ctx paramType:context.Context, paramName:info paramType:struct{}}
func (c CGenService) getFuncParams(node *ast.FuncDecl) (params []map[string]string) {
	if node.Type.Params == nil {
		return
	}
	for _, param := range node.Type.Params.List {
		if param.Names == nil {
			// No name for the return value.
			resultType, err := c.astExprToString(param.Type)
			if err != nil {
				continue
			}
			params = append(params, map[string]string{
				"paramName": "",
				"paramType": resultType,
			})
			continue
		}
		for _, name := range param.Names {
			paramType, err := c.astExprToString(param.Type)
			if err != nil {
				continue
			}
			params = append(params, map[string]string{
				"paramName": name.Name,
				"paramType": paramType,
			})
		}
	}
	return
}

// getFuncResults retrieves the output parameters of the function.
// It returns the name and type of the output parameters.
// For example:
//
// []map[string]string{resultName:list resultType:[]*User, resultName:err resultType:error}
// []map[string]string{resultName: "", resultType: error}
func (c CGenService) getFuncResults(node *ast.FuncDecl) (results []map[string]string) {
	if node.Type.Results == nil {
		return
	}
	for _, result := range node.Type.Results.List {
		if result.Names == nil {
			// No name for the return value.
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"resultName": "",
				"resultType": resultType,
			})
			continue
		}
		for _, name := range result.Names {
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"resultName": name.Name,
				"resultType": resultType,
			})
		}
	}
	return
}

// getFuncComment retrieves the comment of the function.
func (c CGenService) getFuncComment(node *ast.FuncDecl) string {
	return c.astCommentToString(node.Doc)
}