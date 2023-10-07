package automaton

import "errors"

var (
	ErrUndeclaredVariableSymbol               = errors.New("variable could not be found in the environment")
	ErrInterpretationDidNotResultInAValueType = errors.New("the type of the interpretation result was not a value")
	ErrInvalidTypes                           = errors.New("the interpretation of the expression could not result in a value because of a type error")
	ErrUnhandledExpression                    = errors.New("the expression was unhandled")
)

type Interpreter struct {
	Environment Environment
}

func (interpreter *Interpreter) Evaluate(expression Expression) (Value, error) {
	interpretation, err := expression.accept(interpreter)
	value, ok := interpretation.(Value)
	if !ok {
		return value, ErrInterpretationDidNotResultInAValueType
	}

	return value, err
}

func (interpreter *Interpreter) Execute(statement Statement) error {
	return statement.accept(interpreter)
}

func (interpreter *Interpreter) VisitConstantExpression(constant ConstantExpression) (any, error) {
	return constant.Value, nil
}

func (interpreter *Interpreter) VisitUnaryExpression(unary UnaryExpression) (any, error) {
	operand, err := AcceptExpression[Value](interpreter, unary.Expression)
	if err != nil {
		return nil, err
	}

	if unary.Operator == LogicalNegation {
		operandBool, err := operand.AsBoolean()
		if err != nil {
			return nil, err
		}

		return CreateBoolean(!operandBool), nil
	}

	return nil, ErrUnhandledExpression
}

func (interpreter *Interpreter) VisitBinaryExpression(binary BinaryExpression) (any, error) {
	lhs, err := AcceptExpression[Value](interpreter, binary)
	if err != nil {
		return nil, err
	}

	rhs, err := AcceptExpression[Value](interpreter, binary)
	if err != nil {
		return nil, err
	}

	if binary.Operator == LogicalAnd {
		return lhs.LogicalAnd(rhs)
	} else if binary.Operator == LogicalOr {
		return lhs.LogicalOr(rhs)
	}

	return nil, ErrUnhandledExpression
}

func (interpreter *Interpreter) VisitVariable(variable VariableExpression) (any, error) {
	value, exists := interpreter.Environment.Lookup(variable.Symbol)
	if !exists {
		return value, ErrUndeclaredVariableSymbol
	}
	return value, nil
}

func (interpreter *Interpreter) VisitAssignmentStatement(assignment AssignmentStatement) error {
	rhs, err := AcceptExpression[Value](interpreter, assignment.Expression)
	if err != nil {
		return err
	}

	interpreter.Environment.Store(assignment.Variable.Symbol, rhs)

	return nil
}

func (interpreter *Interpreter) VisitBlockStatement(block BlockStatement) error {
	for idx := range block.Statements {
		if err := block.Statements[idx].accept(interpreter); err != nil {
			return err
		}
	}
	return nil
}
