package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment, line int) object.Object {
	switch node := node.(type) {

	//文
	case *ast.Program:
		return evalProgram(node, env, line)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, line)

	case *ast.LetStatement:
		val := Eval(node.Value, env, line)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env, line)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, line)

	case *ast.InfixExpression:
		left := Eval(node.Left, env, line)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env, line)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, line)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env, line)

	case *ast.IfExpression:
		return evalIfExpression(node, env, line)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, line)

		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.Identifier:
		return evalIdentifier(node, env, line)

	case *ast.CallExpression:
		function := Eval(node.Function, env, line)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env, line)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, line)

	//式
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, line)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env, line)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env, line)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, line)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, line)

	case *ast.CommentStatement:
		return nil
	}
	return nil
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, line int) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env, line)

		if result != nil {
			rt := result.Type()
			if rt == object.RETRUN_VALUE_OBJ || rt == object.ERRIE_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, line int) object.Object {
	condition := Eval(ie.Condition, env, line)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env, line)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env, line)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object, line int) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, line)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, line)
	case left.Type() != right.Type():
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s type mismatch: %s %s %s", lineInError, left.Type(), operator, right.Type())
	default:
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unknown operator: %s %s %s", lineInError, left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, line int) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unknown operator: %s %s %s", lineInError, left.Type(), operator, right.Type())
	}
}

func evalProgram(program *ast.Program, env *object.Environment, line int) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env, line)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object, line int) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, line)
	default:
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unknown operator: %s%s", lineInError, operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, line int) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unknown operator: -%s", lineInError, right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERRIE_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment, line int) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	lineInError := fmt.Sprintf("line %d : ", line)
	return newError(lineInError + "identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment, line int) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env, line)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object, line int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv, line)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s not a function: %s", lineInError, fn.Type())
	}

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object, line int) object.Object {
	if operator != "+" {
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unknown operator: %s %s %s", lineInError, left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIndexExpression(left, index object.Object, line int) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index, line)
	default:
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s index operator not supported: %s", lineInError, left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment, line int) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env, line)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			lineInError := fmt.Sprintf("line %d :", line)
			return newError("%s unusable as hash key: %s", lineInError, key.Type())
		}

		value := Eval(valueNode, env, line)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object, line int) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		lineInError := fmt.Sprintf("line %d :", line)
		return newError("%s unusable as hash key: %s", lineInError, index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
