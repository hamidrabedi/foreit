package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
)

// Diagnostic reports a model expression that code generation cannot evaluate.
type Diagnostic struct {
	File    string // absolute or as-parsed path
	Line    int
	Column  int
	Model   string // model type name, empty if unknown
	Method  string // "Fields", "Relations" or "Meta"
	Message string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s:%d:%d: %s.%s: %s", d.File, d.Line, d.Column, d.Model, d.Method, d.Message)
}

// Diagnostics returns parser diagnostics in source order.
func (p *ASTParser) Diagnostics() []Diagnostic {
	diagnostics := append([]Diagnostic(nil), p.diagnostics...)
	sort.Slice(diagnostics, func(i, j int) bool {
		if diagnostics[i].File != diagnostics[j].File {
			return diagnostics[i].File < diagnostics[j].File
		}
		if diagnostics[i].Line != diagnostics[j].Line {
			return diagnostics[i].Line < diagnostics[j].Line
		}
		return diagnostics[i].Column < diagnostics[j].Column
	})
	return diagnostics
}

func (p *ASTParser) reportDiagnostic(pos token.Pos, model, method, message string) {
	position := p.fset.Position(pos)
	p.diagnostics = append(p.diagnostics, Diagnostic{File: position.Filename, Line: position.Line, Column: position.Column, Model: model, Method: method, Message: message})
}

func (p *ASTParser) reportUnsupportedExpressions(model, methodName string, method *ast.FuncDecl) {
	if method.Body == nil {
		return
	}
	assigned := collectAssignedExprs(method.Body)
	seenSlices := make(map[*ast.CompositeLit]bool)
	for _, stmt := range method.Body.List {
		switch stmt := stmt.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			p.reportDiagnostic(stmt.Pos(), model, methodName, "loops cannot be evaluated during generation; list each field explicitly")
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			p.reportDiagnostic(stmt.Pos(), model, methodName, "conditional assembly cannot be evaluated during generation; list each field explicitly")
		}

		for _, expr := range methodExpressions(stmt) {
			if slice, ok := p.resolveAssignedExpr(expr, assigned).(*ast.CompositeLit); ok && !seenSlices[slice] {
				seenSlices[slice] = true
				p.reportUnsupportedSliceElements(slice, assigned, model, methodName)
			}
			if call, ok := expr.(*ast.CallExpr); ok && isAppend(call) {
				p.reportUnsupportedAppend(call, assigned, model, methodName)
			}
		}
	}
}

func methodExpressions(stmt ast.Stmt) []ast.Expr {
	switch stmt := stmt.(type) {
	case *ast.ReturnStmt:
		return stmt.Results
	case *ast.AssignStmt:
		return stmt.Rhs
	case *ast.DeclStmt:
		if decl, ok := stmt.Decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			var expressions []ast.Expr
			for _, spec := range decl.Specs {
				if value, ok := spec.(*ast.ValueSpec); ok {
					expressions = append(expressions, value.Values...)
				}
			}
			return expressions
		}
	case *ast.ExprStmt:
		return []ast.Expr{stmt.X}
	}
	return nil
}

func (p *ASTParser) reportUnsupportedSliceElements(slice *ast.CompositeLit, assigned map[string]ast.Expr, model, methodName string) {
	for _, element := range slice.Elts {
		if name, ok := computedValue(element, assigned); ok {
			p.reportDiagnostic(element.Pos(), model, methodName, fmt.Sprintf("value of %s is computed at run time and cannot be evaluated during generation", name))
			continue
		}
		call, ok := element.(*ast.CallExpr)
		if !ok || p.isSchemaConstructor(methodName, call) {
			continue
		}
		p.reportDiagnostic(element.Pos(), model, methodName, fmt.Sprintf("helper call %s(...) cannot be evaluated during generation; declare fields with direct schema constructor calls", callName(call)))
	}
}

func (p *ASTParser) reportUnsupportedAppend(call *ast.CallExpr, assigned map[string]ast.Expr, model, methodName string) {
	for _, arg := range call.Args[1:] {
		if nested, ok := arg.(*ast.CallExpr); ok && p.isSchemaConstructor(methodName, nested) {
			continue
		}
		if call.Ellipsis.IsValid() || isComputedSlice(arg, assigned) {
			p.reportDiagnostic(arg.Pos(), model, methodName, "append of a computed slice cannot be evaluated during generation")
			continue
		}
		if nested, ok := arg.(*ast.CallExpr); ok && !p.isSchemaConstructor(methodName, nested) {
			p.reportDiagnostic(arg.Pos(), model, methodName, fmt.Sprintf("helper call %s(...) cannot be evaluated during generation; declare fields with direct schema constructor calls", callName(nested)))
		}
	}
}

func computedValue(expr ast.Expr, assigned map[string]ast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		value, ok := assigned[expr.Name]
		if ok {
			if _, ok := value.(*ast.CallExpr); ok {
				return expr.Name, true
			}
		}
	case *ast.SelectorExpr:
		if _, ok := computedValue(expr.X, assigned); ok {
			return expressionName(expr), true
		}
		if _, ok := expr.X.(*ast.CallExpr); ok {
			return expressionName(expr), true
		}
	}
	return "", false
}

func isComputedSlice(expr ast.Expr, assigned map[string]ast.Expr) bool {
	if _, ok := expr.(*ast.CallExpr); ok {
		return true
	}
	if _, ok := computedValue(expr, assigned); ok {
		return true
	}
	_, ok := expr.(*ast.SliceExpr)
	return ok
}

func isAppend(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "append" && len(call.Args) > 1
}

func (p *ASTParser) isSchemaConstructor(methodName string, call *ast.CallExpr) bool {
	switch methodName {
	case "Fields":
		return p.extractFieldFromCall(call) != nil
	case "Relations":
		return p.extractRelationFromExpr(call) != nil
	}
	return false
}

func callName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		return fun.Sel.Name
	}
	return "expression"
}

func expressionName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return expressionName(expr.X) + "." + expr.Sel.Name
	case *ast.CallExpr:
		return callName(expr) + "()"
	}
	return "expression"
}
