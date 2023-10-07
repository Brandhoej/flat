package automaton

import "errors"

var ErrValueHasIncorrectType = errors.New("the value was not a boolean")

type Value struct {
	Type     Type
	Concrete any
}

func CreateBoolean(concrete bool) Value {
	return Value{
		Type:     Boolean,
		Concrete: concrete,
	}
}

func (value Value) IsBoolean() bool {
	return value.Type == Boolean
}

func (value Value) AsBoolean() (bool, error) {
	var boolean bool

	if !value.IsBoolean() {
		return boolean, ErrValueHasIncorrectType
	}

	boolean, ok := value.Concrete.(bool)
	if !ok {
		return boolean, ErrValueHasIncorrectType
	}

	return boolean, nil
}

func (lhs Value) LogicalOr(rhs Value) (Value, error) {
	lhsBool, err := lhs.AsBoolean()
	if err != nil {
		return Value{}, errors.Join(ErrValueHasIncorrectType, err)
	}

	rhsBool, err := rhs.AsBoolean()
	if err != nil {
		return Value{}, errors.Join(ErrValueHasIncorrectType, err)
	}

	return CreateBoolean(lhsBool || rhsBool), nil
}

func (lhs Value) LogicalAnd(rhs Value) (Value, error) {
	lhsBool, err := lhs.AsBoolean()
	if err != nil {
		return Value{}, ErrValueHasIncorrectType
	}

	rhsBool, err := rhs.AsBoolean()
	if err != nil {
		return Value{}, ErrValueHasIncorrectType
	}

	return CreateBoolean(lhsBool && rhsBool), nil
}

func (operand Value) LogicalNegation() (Value, error) {
	operandBool, err := operand.AsBoolean()
	if err != nil {
		return Value{}, errors.Join(ErrValueHasIncorrectType, err)
	}

	return CreateBoolean(!operandBool), nil
}
