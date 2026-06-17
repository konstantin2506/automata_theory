package interpreter

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrLoopCounterNotInt = errors.New("for loop counter type is not int")
	ErrLoopStopperNotInt = errors.New("for loop stopper type is not int")
	ErrLoopStepperNotInt = errors.New("for loop stepper type is not int")
)

type ForStatementNode struct {
	counterName string
	stopperName string
	stepperName string
	doThis      AstNode
}

func NewForStatementNode(counter, stopper, stepper string, doThis AstNode) *ForStatementNode {
	return &ForStatementNode{counter, stopper, stepper, doThis}
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

func (node *ForStatementNode) Eval(scope *Scope) (Variable, error) {
	counterFound, err := scope.FindVariableDepth(node.counterName)
	if err != nil {
		return nil, err
	}
	stepperFound, err := scope.FindVariableDepth(node.stepperName)
	if err != nil {
		return nil, err
	}
	stopperFound, err := scope.FindVariableDepth(node.stopperName)
	if err != nil {
		return nil, err
	}

	ok := counterFound.Type() == Array && stopperFound.Type() == Array && stepperFound.Type() == Array
	if ok && counterFound.(*VarArray).InnerType() == Int && stepperFound.(*VarArray).InnerType() == Int && stopperFound.(*VarArray).InnerType() == Int {
		c := counterFound.(*VarArray)
		ste := stepperFound.(*VarArray)
		sto := stopperFound.(*VarArray)
		if slices.Equal(c.sizes, ste.sizes) && slices.Equal(c.sizes, sto.sizes) {
			counterData := 0
			stopperData := 0
			stepperData := 0
			for i := range c.data {
				counterData = c.data[i].(*Integer).Data()
				stopperData = sto.data[i].(*Integer).Data()
				stepperData = ste.data[i].(*Integer).Data()
				predicate := func(x, y int) bool { return x < y }
				if stepperData < 0 {
					predicate = func(x, y int) bool { return x > y }
				}
				for i := counterData; predicate(i, stopperData); i += stepperData {
					c.data[i].(*Integer).data = i
					childScope := NewScope(scope, scope.globalScope)
					res, err := node.doThis.Eval(&childScope)
					if err != nil {
						return nil, err
					}
					if res != nil {
						err := scope.AssignScalar(node.counterName, NewVariableInt(i*stepperData))
						if err != nil {
							return nil, err
						}
						return res, nil
					}

				}

			}
		} else {
			return nil, ErrVarInvalidType
		}
	} else {

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
		predicate := func(x, y int) bool { return x < y }
		if stepper.Data() < 0 {
			predicate = func(x, y int) bool { return x > y }
		}
		for i := counter.Data(); predicate(i, stopper.Data()); i += stepper.Data() {
			counter.data = i
			childScope := NewScope(scope, scope.globalScope)
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
		// err = scope.AssignScalar(node.counterName, NewVariableInt(stopper.Data()*stepper.Data()))

		return nil, nil
	}
	return nil, nil
}
