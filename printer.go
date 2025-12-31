package main

// type AstPrinter struct {
// }
//
// func (ap *AstPrinter) print(expr Expr) any {
// 	return expr.accept(ap)
// }
//
// func (ap *AstPrinter) visitBinaryExpr(expr *Binary) any {
// 	return ap.parenthesize(expr.operator.lexeme, expr.left, expr.right)
// }
//
// func (ap *AstPrinter) visitGroupingExpr(expr *Grouping) any {
// 	return ap.parenthesize("group", expr.expression)
// }
//
// func (ap *AstPrinter) visitLiteralExpr(expr *Literal) any {
// 	if expr.value == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprint(expr.value)
// }
//
// func (ap *AstPrinter) visitUnaryExpr(expr *Unary) any {
// 	return ap.parenthesize(expr.operator.lexeme, expr.right)
// }
//
// func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
// 	var builder = ""
// 	builder += "(" + name
// 	for _, expr := range exprs {
// 		if expr != nil {
// 			builder += " " + fmt.Sprint(expr.accept(ap))
// 		}
// 	}
// 	builder += ")"
// 	return builder
// }
