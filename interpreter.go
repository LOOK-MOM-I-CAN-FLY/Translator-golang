package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"translator/parser"
)

// Value represents a runtime value.
type Value interface {
	Type() string
	String() string
}

// ==== Primitive values ====

type IntValue int64

func (v IntValue) Type() string   { return "int" }
func (v IntValue) String() string { return fmt.Sprintf("%d", v) }

type StringValue string

func (v StringValue) Type() string   { return "string" }
func (v StringValue) String() string { return string(v) }

type BoolValue bool

func (v BoolValue) Type() string   { return "bool" }
func (v BoolValue) String() string { return fmt.Sprintf("%v", bool(v)) }

// ==== Struct support ====

type TypeDef struct {
	Name   string
	Fields map[string]string
}

type Instance struct {
	Def    *TypeDef
	Fields map[string]Value
}

func (i *Instance) Type() string {
	if i.Def != nil {
		return i.Def.Name
	}
	return "object"
}

func (i *Instance) String() string {
	if i.Def == nil {
		return fmt.Sprintf("%v", i.Fields)
	}
	parts := make([]string, 0, len(i.Fields))
	for k, v := range i.Fields {
		parts = append(parts, fmt.Sprintf("%s:%s", k, v.String()))
	}
	return fmt.Sprintf("%s{%s}", i.Def.Name, strings.Join(parts, ", "))
}

// ==== Environment with scopes ====

type Environment struct {
	scopes []map[string]Value
}

func NewEnvironment() *Environment {
	return &Environment{scopes: []map[string]Value{make(map[string]Value)}}
}

func (e *Environment) Push() {
	e.scopes = append(e.scopes, make(map[string]Value))
}

func (e *Environment) Pop() {
	if len(e.scopes) > 1 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *Environment) Declare(name string, value Value) {
	e.scopes[len(e.scopes)-1][name] = value
}

func (e *Environment) DeclareOrSet(name string, value Value) {
	if _, ok := e.scopes[len(e.scopes)-1][name]; ok {
		e.scopes[len(e.scopes)-1][name] = value
		return
	}
	e.Declare(name, value)
}

func (e *Environment) Set(name string, value Value) error {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if _, ok := e.scopes[i][name]; ok {
			e.scopes[i][name] = value
			return nil
		}
	}
	return fmt.Errorf("undefined variable: %s", name)
}

func (e *Environment) Get(name string) (Value, error) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		if v, ok := e.scopes[i][name]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("undefined variable: %s", name)
}

// ==== Interpreter ====

type Interpreter struct {
	*parser.BaseSimpleParserVisitor
	env    *Environment
	types  map[string]*TypeDef
	output []string
	err    error
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		BaseSimpleParserVisitor: &parser.BaseSimpleParserVisitor{},
		env:                     NewEnvironment(),
		types:                   make(map[string]*TypeDef),
		output:                  []string{},
	}
}

func (i *Interpreter) Err() error {
	return i.err
}

func (i *Interpreter) Visit(tree antlr.ParseTree) interface{} {
	if tree == nil {
		return nil
	}
	return tree.Accept(i)
}

func (i *Interpreter) failf(format string, args ...interface{}) {
	if i.err == nil {
		i.err = fmt.Errorf(format, args...)
	}
}

func (i *Interpreter) VisitChildren(node antlr.RuleNode) interface{} {
	if i.err != nil {
		return nil
	}
	var result interface{}
	for j := 0; j < node.GetChildCount(); j++ {
		child := node.GetChild(j)
		if tree, ok := child.(antlr.ParseTree); ok {
			result = tree.Accept(i)
			if i.err != nil {
				return nil
			}
		}
	}
	return result
}

// program
func (i *Interpreter) VisitProgram(ctx *parser.ProgramContext) interface{} {
	return i.VisitChildren(ctx)
}

func (i *Interpreter) VisitTypeDeclaration(ctx *parser.TypeDeclarationContext) interface{} {
	return i.VisitChildren(ctx)
}

func (i *Interpreter) VisitStatement(ctx *parser.StatementContext) interface{} {
	return i.VisitChildren(ctx)
}

// type declarations
func (i *Interpreter) VisitStructDeclaration(ctx *parser.StructDeclarationContext) interface{} {
	if i.err != nil {
		return nil
	}
	name := ctx.IDENTIFIER().GetText()
	fields := make(map[string]string)
	for _, f := range ctx.AllStructField() {
		fieldName := f.IDENTIFIER().GetText()
		fieldType := f.Type_().GetText()
		fields[fieldName] = fieldType
	}
	i.types[name] = &TypeDef{Name: name, Fields: fields}
	return nil
}

// variable declaration
func (i *Interpreter) VisitDeclaration(ctx *parser.DeclarationContext) interface{} {
	if i.err != nil {
		return nil
	}
	name := ctx.IDENTIFIER().GetText()
	varType := ""
	if ctx.Type_() != nil {
		varType = ctx.Type_().GetText()
	}

	var value Value
	if ctx.Expression() != nil {
		value = i.evalExpression(ctx.Expression())
		if i.err != nil {
			return nil
		}
	}

	if varType == "" {
		if value == nil {
			i.failf("variable %s must have a type or value", name)
			return nil
		}
		varType = value.Type()
	}

	if value == nil {
		value = i.defaultValue(varType)
		if i.err != nil {
			return nil
		}
	} else {
		value = i.ensureType(value, varType)
		if i.err != nil {
			return nil
		}
	}

	i.env.Declare(name, value)
	return nil
}

// assignment
func (i *Interpreter) VisitAssignment(ctx *parser.AssignmentContext) interface{} {
	if i.err != nil {
		return nil
	}
	// short declaration
	if ctx.DECLARE() != nil {
		name := ctx.IDENTIFIER().GetText()
		val := i.evalExpression(ctx.Expression())
		if i.err != nil {
			return nil
		}
		i.env.DeclareOrSet(name, val)
		return nil
	}

	// regular assignment to lvalue
	val := i.evalExpression(ctx.Expression())
	if i.err != nil {
		return nil
	}
	lv := ctx.Lvalue()
	if lv == nil {
		i.failf("expected lvalue for assignment")
		return nil
	}
	i.assignLValue(lv.(*parser.LvalueContext), val)
	return nil
}

// if statement
func (i *Interpreter) VisitIfStatement(ctx *parser.IfStatementContext) interface{} {
	if i.err != nil {
		return nil
	}
	cond := i.evalExpression(ctx.Expression())
	if i.err != nil {
		return nil
	}
	blocks := ctx.AllBlock()
	if i.toBool(cond) {
		blocks[0].Accept(i)
		return nil
	}
	if len(blocks) > 1 {
		blocks[1].Accept(i)
	}
	return nil
}

// for statement
func (i *Interpreter) VisitForStatement(ctx *parser.ForStatementContext) interface{} {
	if i.err != nil {
		return nil
	}
	body := ctx.Block()

	if fc := ctx.ForClause(); fc != nil {
		clause := fc.(*parser.ForClauseContext)
		if clause.GetInit() != nil {
			clause.GetInit().(antlr.ParseTree).Accept(i)
			if i.err != nil {
				return nil
			}
		}

		for {
			if clause.GetCond() != nil {
				condVal := i.evalExpression(clause.GetCond())
				if i.err != nil {
					return nil
				}
				if !i.toBool(condVal) {
					break
				}
			}

			body.Accept(i)
			if i.err != nil {
				return nil
			}

			if clause.GetPost() != nil {
				clause.GetPost().(antlr.ParseTree).Accept(i)
				if i.err != nil {
					return nil
				}
			}
		}
		return nil
	}

	if ctx.Expression() != nil {
		for {
			condVal := i.evalExpression(ctx.Expression())
			if i.err != nil {
				return nil
			}
			if !i.toBool(condVal) {
				break
			}
			body.Accept(i)
			if i.err != nil {
				return nil
			}
		}
		return nil
	}

	// for { ... } - protect from infinite loop without break
	maxIterations := 100000
	for iter := 0; iter < maxIterations; iter++ {
		body.Accept(i)
		if i.err != nil {
			return nil
		}
	}
	i.failf("for { } exceeded %d iterations without break", maxIterations)
	return nil
}

// block with scope
func (i *Interpreter) VisitBlock(ctx *parser.BlockContext) interface{} {
	if i.err != nil {
		return nil
	}
	i.env.Push()
	defer i.env.Pop()
	for _, st := range ctx.AllStatement() {
		st.Accept(i)
		if i.err != nil {
			return nil
		}
	}
	return nil
}

// function call: fmt.Println
func (i *Interpreter) VisitFunctionCall(ctx *parser.FunctionCallContext) interface{} {
	if i.err != nil {
		return nil
	}
	args := []string{}
	if ctx.ArgList() != nil {
		for _, expr := range ctx.ArgList().AllExpression() {
			val := i.evalExpression(expr)
			if i.err != nil {
				return nil
			}
			args = append(args, val.String())
		}
	}
	out := strings.Join(args, " ")
	fmt.Println(out)
	i.output = append(i.output, out)
	return nil
}

// ==== Expressions ====

func (i *Interpreter) evalExpression(ctx parser.IExpressionContext) Value {
	val := i.Visit(ctx)
	if i.err != nil {
		return nil
	}
	if v, ok := val.(Value); ok {
		return v
	}
	i.failf("expression did not produce a value")
	return nil
}

func (i *Interpreter) VisitExpression(ctx *parser.ExpressionContext) interface{} {
	return i.Visit(ctx.OrExpr())
}

func (i *Interpreter) VisitOrExpr(ctx *parser.OrExprContext) interface{} {
	left := i.Visit(ctx.AndExpr(0)).(Value)
	for idx := 1; idx < len(ctx.AllAndExpr()); idx++ {
		right := i.Visit(ctx.AndExpr(idx)).(Value)
		left = BoolValue(i.toBool(left) || i.toBool(right))
	}
	return left
}

func (i *Interpreter) VisitAndExpr(ctx *parser.AndExprContext) interface{} {
	left := i.Visit(ctx.CompExpr(0)).(Value)
	for idx := 1; idx < len(ctx.AllCompExpr()); idx++ {
		right := i.Visit(ctx.CompExpr(idx)).(Value)
		left = BoolValue(i.toBool(left) && i.toBool(right))
	}
	return left
}

func (i *Interpreter) VisitCompExpr(ctx *parser.CompExprContext) interface{} {
	left := i.Visit(ctx.AddExpr(0)).(Value)
	ops := ctx.AllCompOp()
	for idx, opNode := range ops {
		right := i.Visit(ctx.AddExpr(idx + 1)).(Value)
		left = BoolValue(i.compare(left, opNode.GetText(), right))
	}
	return left
}

func (i *Interpreter) VisitAddExpr(ctx *parser.AddExprContext) interface{} {
	left := i.Visit(ctx.MulExpr(0)).(Value)
	ops := ctx.AllAddOp()
	for idx, opNode := range ops {
		right := i.Visit(ctx.MulExpr(idx + 1)).(Value)
		op := opNode.GetText()
		if op == "+" {
			if left.Type() == "string" || right.Type() == "string" {
				left = StringValue(left.String() + right.String())
			} else {
				left = IntValue(i.toInt(left) + i.toInt(right))
			}
		} else {
			left = IntValue(i.toInt(left) - i.toInt(right))
		}
	}
	return left
}

func (i *Interpreter) VisitMulExpr(ctx *parser.MulExprContext) interface{} {
	left := i.Visit(ctx.UnaryExpr(0)).(Value)
	ops := ctx.AllMulOp()
	for idx, opNode := range ops {
		right := i.Visit(ctx.UnaryExpr(idx + 1)).(Value)
		op := opNode.GetText()
		if op == "*" {
			left = IntValue(i.toInt(left) * i.toInt(right))
		} else {
			left = IntValue(i.toInt(left) / i.toInt(right))
		}
	}
	return left
}

func (i *Interpreter) VisitUnaryExpr(ctx *parser.UnaryExprContext) interface{} {
	if ctx.Primary() != nil {
		return i.Visit(ctx.Primary())
	}
	op := ctx.GetChild(0).(antlr.ParseTree).GetText()
	val := i.Visit(ctx.UnaryExpr()).(Value)
	if op == "!" {
		return BoolValue(!i.toBool(val))
	}
	return IntValue(-i.toInt(val))
}

func (i *Interpreter) VisitPrimary(ctx *parser.PrimaryContext) interface{} {
	if ctx.Literal() != nil {
		return i.Visit(ctx.Literal())
	}
	if ctx.Expression() != nil {
		return i.Visit(ctx.Expression())
	}

	// identifier (with optional field access)
	names := ctx.AllIDENTIFIER()
	if len(names) == 0 {
		i.failf("expected primary expression")
		return nil
	}

	baseName := names[0].GetText()
	val, err := i.env.Get(baseName)
	if err != nil {
		i.failf(err.Error())
		return nil
	}

	for idx := 1; idx < len(names); idx++ {
		field := names[idx].GetText()
		inst, ok := val.(*Instance)
		if !ok {
			i.failf("value of type %s has no fields", val.Type())
			return nil
		}
		fieldVal, ok := inst.Fields[field]
		if !ok {
			i.failf("unknown field '%s' in %s", field, inst.Type())
			return nil
		}
		val = fieldVal
	}

	return val
}

func (i *Interpreter) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	if ctx.INT_LIT() != nil {
		v, _ := strconv.ParseInt(ctx.INT_LIT().GetText(), 10, 64)
		return IntValue(v)
	}
	if ctx.STRING_LIT() != nil {
		text := ctx.STRING_LIT().GetText()
		if len(text) >= 2 && text[0] == '"' {
			text = text[1 : len(text)-1]
		}
		return StringValue(text)
	}
	if ctx.TRUE() != nil {
		return BoolValue(true)
	}
	return BoolValue(false)
}

// ==== Helpers ====

func (i *Interpreter) defaultValue(typeName string) Value {
	switch typeName {
	case "int":
		return IntValue(0)
	case "string":
		return StringValue("")
	case "bool":
		return BoolValue(false)
	default:
		def, ok := i.types[typeName]
		if !ok {
			i.failf("unknown type: %s", typeName)
			return nil
		}
		return i.newInstance(def)
	}
}

func (i *Interpreter) newInstance(def *TypeDef) *Instance {
	fields := make(map[string]Value)
	for name, t := range def.Fields {
		fields[name] = i.defaultValue(t)
		if i.err != nil {
			return nil
		}
	}
	return &Instance{Def: def, Fields: fields}
}

func (i *Interpreter) ensureType(val Value, typeName string) Value {
	if val == nil {
		return nil
	}
	if val.Type() == typeName {
		return val
	}
	// allow assignment of instance to its type name
	if inst, ok := val.(*Instance); ok && inst.Type() == typeName {
		return val
	}
	i.failf("cannot assign %s to %s", val.Type(), typeName)
	return nil
}

func (i *Interpreter) assignLValue(ctx *parser.LvalueContext, val Value) {
	names := ctx.AllIDENTIFIER()
	if len(names) == 0 {
		i.failf("invalid lvalue")
		return
	}

	if len(names) == 1 {
		name := names[0].GetText()
		if err := i.env.Set(name, val); err != nil {
			i.failf(err.Error())
		}
		return
	}

	baseName := names[0].GetText()
	obj, err := i.env.Get(baseName)
	if err != nil {
		i.failf(err.Error())
		return
	}

	for idx := 1; idx < len(names)-1; idx++ {
		field := names[idx].GetText()
		inst, ok := obj.(*Instance)
		if !ok {
			i.failf("value of type %s has no fields", obj.Type())
			return
		}
		fieldVal, ok := inst.Fields[field]
		if !ok {
			i.failf("unknown field '%s' in %s", field, inst.Type())
			return
		}
		obj = fieldVal
	}

	lastField := names[len(names)-1].GetText()
	inst, ok := obj.(*Instance)
	if !ok {
		i.failf("value of type %s has no fields", obj.Type())
		return
	}
	if _, ok := inst.Fields[lastField]; !ok {
		i.failf("unknown field '%s' in %s", lastField, inst.Type())
		return
	}
	inst.Fields[lastField] = val
}

func (i *Interpreter) toInt(v Value) int64 {
	switch val := v.(type) {
	case IntValue:
		return int64(val)
	case BoolValue:
		if val {
			return 1
		}
		return 0
	case StringValue:
		parsed, err := strconv.ParseInt(string(val), 10, 64)
		if err == nil {
			return parsed
		}
	}
	return 0
}

func (i *Interpreter) toBool(v Value) bool {
	switch val := v.(type) {
	case BoolValue:
		return bool(val)
	case IntValue:
		return val != 0
	case StringValue:
		return len(val) > 0
	}
	return false
}

func (i *Interpreter) compare(left Value, op string, right Value) bool {
	switch op {
	case "==":
		return left.String() == right.String()
	case "!=":
		return left.String() != right.String()
	case "<":
		return i.toInt(left) < i.toInt(right)
	case ">":
		return i.toInt(left) > i.toInt(right)
	case "<=":
		return i.toInt(left) <= i.toInt(right)
	case ">=":
		return i.toInt(left) >= i.toInt(right)
	}
	return false
}
