package automaton

type StatementVisitor interface {
	VisitAssignmentStatement(assignment AssignmentStatement) error
	VisitBlockStatement(block BlockStatement) error
}

type Statement interface {
	accept(visitor StatementVisitor) error
}

type AssignmentStatement struct {
	Variable   VariableExpression
	Expression Expression
}

func (assignment AssignmentStatement) accept(visitor StatementVisitor) error {
	return visitor.VisitAssignmentStatement(assignment)
}

type BlockStatement struct {
	Statements []Statement
}

func (block BlockStatement) accept(visitor StatementVisitor) error {
	return visitor.VisitBlockStatement(block)
}
