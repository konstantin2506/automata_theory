package interpreter

import (
	"errors"
	"fmt"
)

var (
	ErrLoopCounterNotInt = errors.New("for loop counter type is not int")
	ErrLoopStopperNotInt = errors.New("for loop stopper type is not int")
	ErrLoopStepperNotInt = errors.New("for loop stepper type is not int")
)

type ForStatement struct {
	counterName string
	stopperName string
	stepperName string
	doThis      AstNode
}

func NewForStatement(counter, stopper, stepper string, doThis AstNode) *ForStatement {
	return &ForStatement{counter, stopper, stepper, doThis}
}

func findInt(scope *Scope, targetName string, errType error) (*Integer, error) {
	target, err := scope.FindVariableDepth(targetName)
	if err != nil {
		return nil, err
	}
	if target.Type() != Int {
		return nil, fmt.Errorf("%w: (name=%s, type=%s)", errType, targetName, TypeName(target.Type()))
	}
	return target.(*Integer), nil
}

func (node *ForStatement) Eval(scope *Scope) (Variable, error) {
	counter, err := findInt(scope, node.counterName, ErrLoopCounterNotInt)
	if err != nil {
		return nil, err
	}
	stopper, err := findInt(scope, node.stopperName, ErrLoopStopperNotInt)
	if err != nil {
		return nil, err
	}
	stepper, err := findInt(scope, node.stepperName, ErrLoopStepperNotInt)
	if err != nil {
		return nil, err
	}
	for i := counter.Data(); i < stopper.Data(); i += stepper.Data() {
		childScope := NewScope(scope)
		res, err := node.doThis.Eval(&childScope)
		if err != nil {
			return nil, err
		}
		if res != nil {
			err := scope.AssignScalar(node.counterName, NewVariableInt(i*stepper.Data()))
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	err = scope.AssignScalar(node.counterName, NewVariableInt(stopper.Data()*stepper.Data()))

	return nil, err
}
