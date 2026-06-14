package interpreter

import (
	"errors"
	"fmt"
)

var (
	ErrArraySizes             = errors.New("array sizes are incorrect")
	ErrArrayIndices           = errors.New("array indices count differ")
	ErrArrayOutOfRange        = errors.New("array index is out of range")
	ErrArrayAssignTypesDiffer = errors.New("array assign types differ")
	ErrArrayDimsDiffer        = errors.New("array dims differ")
)

type VarArray struct {
	sizes []int
	data  []Variable
	dataT VarT
}

func CmpTypeWithInner(vector Vector, value Variable) bool {
	return vector.InnerType() == value.Type()
}

func (v *VarArray) InnerType() VarT {
	return v.dataT
}

func (v *VarArray) CmpTypeWith(other *VarArray) error {
	if !CmpTypeWithInner(v, other.data[0]) {
		return fmt.Errorf("%w: want: %d, got: %d", ErrArrayAssignTypesDiffer, v.dataT, other.InnerType())
	}
	if len(v.sizes) != len(other.sizes) {
		return fmt.Errorf("%w: sizes: %d, other: %d", ErrArrayIndices, len(v.sizes), len(other.sizes))
	}
	for j, size := range v.sizes {
		if size != other.sizes[j] {
			return ErrArrayDimsDiffer
		}
	}
	return nil
}

func (v *VarArray) At(indices []int) (Variable, int, error) {
	if len(indices) != len(v.sizes) {
		return nil, 0, fmt.Errorf("%w: sizes: %d, indices: %d", ErrArrayIndices, len(v.sizes), len(indices))
	}
	for i := range indices { // indexing from 1
		indices[i] = indices[i] - 1
	}

	resultIndex := indices[len(indices)-1]
	if resultIndex > v.sizes[len(v.sizes)-1] {
		return nil, 0, fmt.Errorf("%w: idx[%d] = %d, dim[%d] = %d", ErrArrayOutOfRange, len(indices)-1, resultIndex, len(indices)-1, v.sizes[len(indices)-1])
	}
	for i := len(indices) - 1; i >= 0; i-- {
		if v.sizes[i] < indices[i] {
			return nil, 0, fmt.Errorf("%w: idx[%d] = %d, dim[%d] = %d", ErrArrayOutOfRange, i, indices[i], i, v.sizes[i])
		}
		resultIndex = resultIndex*v.sizes[i] + indices[i]
	}
	return v.data[resultIndex], resultIndex, nil
}

func (v *VarArray) Assign(indices []int, value Variable) error {
	if !CmpTypeWithInner(v, value) {
		return fmt.Errorf("%w: want: %d, got: %d", ErrArrayAssignTypesDiffer, v.dataT, value.Type())
	}
	_, resultIndex, err := v.At(indices)
	if err != nil {
		return err
	}
	v.data[resultIndex] = value
	return nil
}

func (v *VarArray) Type() VarT {
	return Array
}

func NewArray(sizes []int, value Variable) (Variable, error) {
	if sizes[0] == 0 {
		return nil, fmt.Errorf("%w: first size of array is 0", ErrArraySizes)
	}
	dataSize := 1
	for _, size := range sizes {
		dataSize *= size
		if size == 0 {
			return nil, fmt.Errorf("%w: zero size in middle of sizes", ErrArraySizes)
		}

	}
	sizesCopy := make([]int, len(sizes))
	copy(sizesCopy, sizes)
	data := make([]Variable, dataSize)
	for i := range dataSize {
		data[i] = value.Copy()
	}
	v := &VarArray{sizesCopy, data, value.Type()}
	return v, nil
}

func (v *VarArray) Copy() Variable {
	newData := make([]Variable, len(v.data))
	newSizes := make([]int, len(v.sizes))
	for i := range len(newData) {
		newData[i] = v.data[i].Copy()
	}
	copy(newSizes, v.sizes)
	return &VarArray{newSizes, newData, v.dataT}
}

type ArrayElemNode struct {
	indices []int
	array   AstNode
}

func NewArrayElemNode(indices []int, array AstNode) AstNode {
	return &ArrayElemNode{indices, array}
}

func (node *ArrayElemNode) Eval(scope *Scope) (Variable, error) {
	arr, err := node.array.Eval(scope)
	if err != nil {
		return nil, err
	}
	res, _, err := arr.(*VarArray).At(node.indices)
	return res, err
}

type ArrayDeclNode struct {
	innerT VarT
	sizes  []int
	value  AstNode
}

func NewArrayDeclNode(innerType VarT, sizes []int, value AstNode) AstNode {
	return &ArrayDeclNode{innerType, sizes, value}
}

func (node *ArrayDeclNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.value.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() != node.innerT {
		return nil, fmt.Errorf("%w: want: %s, got: %s", ErrArrayAssignTypesDiffer, TypeName(node.innerT), TypeName(v.Type()))
	}
	return NewArray(node.sizes, v)
}
