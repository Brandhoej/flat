package automaton

import "errors"

var (
	ErrVisitorCouldNotAcceptExpression = errors.New("the visitor resulted in an error when accepting the expression")
	ErrVisitorResultHasIncorrectType   = errors.New("the visitor resulted in an incorrect type being returned without")
)

type ExpressionVisitor interface {
	VisitConstantExpression(constantExpression ConstantExpression) (any, error)
	VisitUnaryExpression(unaryExpression UnaryExpression) (any, error)
	VisitBinaryExpression(binaryExpression BinaryExpression) (any, error)
	VisitVariable(variable VariableExpression) (any, error)
}

type Expression interface {
	accept(visitor ExpressionVisitor) (any, error)
}

func AcceptExpression[T any](visitor ExpressionVisitor, expression Expression) (T, error) {
	var result T
	value, err := expression.accept(visitor)
	if err != nil {
		return result, errors.Join(ErrVisitorCouldNotAcceptExpression, err)
	}

	result, ok := value.(T)
	if !ok {
		return result, ErrVisitorResultHasIncorrectType
	}

	return result, nil
}

type ConstantExpression struct {
	Value Value
}

func (constant ConstantExpression) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitConstantExpression(constant)
}

type VariableExpression struct {
	Symbol Symbol
}

func (variable VariableExpression) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitVariable(variable)
}

type UnaryOperator int32

const (
	LogicalNegation = UnaryOperator(iota)
)

type UnaryExpression struct {
	Operator   UnaryOperator
	Expression Expression
}

func (unaryExpression UnaryExpression) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitUnaryExpression(unaryExpression)
}

type BinaryOperator int32

const (
	LogicalAnd = BinaryOperator(iota)
	LogicalOr
)

type BinaryExpression struct {
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

func (binaryExpression BinaryExpression) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitBinaryExpression(binaryExpression)
}
