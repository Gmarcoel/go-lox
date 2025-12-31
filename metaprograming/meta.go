package main

import (
	"fmt"
	"os"
	"strings"
)

type GenerateAst struct {
	outputDir string
}

func main() {
	var g = GenerateAst{outputDir: "outputs"}
	g.generate([]string{"."})
}

func (gen *GenerateAst) generate(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	gen.outputDir = args[0]
	var expr_list = []string{
		"Assign : Token name, Expr value",
		"Binary : Expr left, Token operator, Expr right",
		"Call : Expr callee, Token paren, []any arguments",
		"Get : Expr object, Token name",
		"Grouping : Expr expression",
		"Literal : any value",
		"Logical : Expr left, Token operator, Expr right",
		"Set : Expr object, Token name, Expr value",
		"Super : Token keyword, Token method",
		"This : Token keyword",
		"Unary : Token operator, Expr right",
		"Variable : Token name",
	}
	var stmt_list = []string{
		"Block : []Stmt statements ",
		"Class : Token name, *Variable superclass, []*Function methods",
		"Expression : Expr expression",
		"Function : Token name, []Token params, []Stmt body",
		"If : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print : Expr expression",
		"Return : Token keyword, Expr value",
		"Va : Token name, Expr initializer",
		"While : Expr condition, Stmt body",
	}
	defineAst(gen.outputDir, "Stmt", stmt_list)

	defineAst(gen.outputDir, "Expr", expr_list)
}

func defineAst(outputDir string, baseName string, types []string) {
	var path = outputDir + "/" + strings.ToLower(baseName) + ".go"
	var writer, _ = os.Create(path)
	writer.WriteString("package main" + "\n")
	writer.WriteString("" + "\n")

	// The base class ()
	writer.WriteString("type " + baseName + " interface {" + "\n")
	writer.WriteString("accept(" + strings.ToLower(baseName) + "Visitor) any" + "\n")
	writer.WriteString("}" + "\n")
	writer.WriteString("" + "\n")
	defineVisitor(writer, baseName, types)

	writer.WriteString("" + "\n")

	// // The AST classes.
	for _, t := range types {
		var className string = strings.TrimSpace(strings.Split(t, ":")[0])
		var fields = strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(writer, baseName, className, fields)
	}

	// writer.WriteString("" + "\n")
	// writer.WriteString("func (" + strings.ToLower(baseName) + " *" + baseName +  ")" + "accept(visitor any) {}" + "\n")
	writer.WriteString("" + "\n")

	writer.Sync()

	return}

type Visitor interface {
	method() string
	method2() int
}

func defineVisitor(writer *os.File, baseName string, types []string) {
	writer.WriteString("type " + strings.ToLower(baseName) + "Visitor interface {" + "\n")
	for _, t := range types {
		var typeName string = strings.TrimSpace(strings.Split(t, ":")[0])
		writer.WriteString("visit" + typeName + baseName + "(" + strings.ToLower(baseName) + " *" + typeName + ") any" + "\n")
	}
	// writer.println(" }");
	writer.WriteString(" }" + "\n")
}

// private static void defineType(
// PrintWriter writer, String baseName,
// String className, String fieldList) {
type a struct {
	b string
	c int
}

func defineType(writer *os.File, baseName string, className string, fieldList string) {
	var fields = strings.Split(fieldList, ", ")
	var fixedFiedList = ""

	// Struct
	writer.WriteString("type " + className + " struct {" + "\n")
	for _, field := range fields {
		var parts = strings.Split(field, " ")
		writer.WriteString(parts[1] + " " + parts[0] + "\n")
		fixedFiedList = fixedFiedList + parts[1] + " " + parts[0] + ", "
	}
	writer.WriteString("}" + "\n")
	writer.WriteString("" + "\n")

	// Visitor pattern.
	writer.WriteString("func (" + strings.ToLower(className) + "_ *" + className + ") accept(visitor " + strings.ToLower(baseName) + "Visitor) any {" + "\n")
	writer.WriteString("return visitor.visit" + className + baseName + "(" + strings.ToLower(className) + "_)" + "\n")
	writer.WriteString("}" + "\n")
	writer.WriteString("" + "\n")

	// Constructor.
	writer.WriteString("func new" + className + "(" + fixedFiedList + ") *" + className + " {" + "\n")
	writer.WriteString("	return &" + className + "{" + "\n")
	for _, field := range fields {
		var name = strings.Split(field, " ")[1]
		writer.WriteString(name + ": " + name + "," + "\n")
	}
	writer.WriteString(" }" + "\n")
	writer.WriteString(" }" + "\n")
}
