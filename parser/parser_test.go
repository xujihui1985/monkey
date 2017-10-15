// Package parser provides
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	t.Skip()
	input := `
  let x = 5; let y = 10;
  let foobar = 100;
  
  `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("should container 3 statements got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
	fmt.Println(program.String())
}

func TestReturnStatment(t *testing.T) {

	t.Skip()
	input := `
  return 5;
  return 10;
  return 1 + 1;
  `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatalf("program statmens does not container 3 satement")
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expect return statement but got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt not 'return' got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {

	t.Skip()
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		t.Log(p.Errors())
		if len(program.Statements) != 1 {
			t.Fatalf("Program statements does not container %d Statements but got %d\n",
				1, len(program.Statements))
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	t.Skip()
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		//{"5-5;", 5, "-", 5},
		//{"5*5;", 5, "*", 5},
		//{"5/5;", 5, "/", 5},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("program statements does not contain %d statements, got %d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement")
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp is not infix expression")
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp opertaor is not %s", tt.operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		//{"-a * b", "((-a) * b)"},
		{"a + b * c", "(a + (b * c))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		actual := program.String()
		if actual != tt.expect {
			t.Fatalf("expect %s acutal %s", tt.expect, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expect %d Statements, but %s", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement")
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expect IfExpression but %T", stmt.Expression)
	}

	t.Log(exp.String())
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("expect %d Statements, but %s", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement")
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral")
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("parameters should be 2")
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams:[]string{"x"}},
		{input: "fn(x, y, z){};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("expected %d params but %d", len(tt.expectedParams), len(function.Parameters))
		}

	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program statements should be 1")
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[0] is not Expression statement")
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not CallExpression")
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("arguments shoud be 3, actual is %d", len(exp.Arguments))
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral should be let but got %q", s.TokenLiteral)
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not LetStatement")
		return false
	}
	if letStmt.Name.Value != name {
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		return false
	}
	return true
}
