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

func (v *VarArray) Print() {
	fmt.Printf("%s_arr", TypeName(v.dataT))
	for _, size := range v.sizes {
		fmt.Printf("[%d]", size)
	}
	fmt.Printf("={")
	for i, elem := range v.data {
		elem.(Printer).Print()
		if i != len(v.data)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}")
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

	resultIndex := 0
	factor := 1
	for i := len(indices) - 1; i >= 0; i-- {
		if v.sizes[i] < indices[i] || indices[i] < 0 {
			return nil, 0, fmt.Errorf("%w: idx[%d] = %d, dim[%d] = %d", ErrArrayOutOfRange, i, indices[i], i, v.sizes[i])
		}
		resultIndex += indices[i] * factor
		factor *= v.sizes[i]
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

func NewArray(sizes []int, value Variable) (*VarArray, error) {
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

func convertSizes(gotSizes []AstNode, scope *Scope) ([]int, error) {
	sizes := make([]int, len(gotSizes))
	for i, sizeNode := range gotSizes {
		s, err := sizeNode.Eval(scope)
		if err != nil {
			return nil, err
		}
		if s.Type() != Int {
			return nil, ErrSizeNotInt
		}
		sizes[i] = s.(*Integer).Data()
	}
	return sizes, nil
}

type ArrayElemNode struct {
	indices []AstNode
	name    string
}

func NewArrayElemNode(indices []AstNode, name string) *ArrayElemNode {
	return &ArrayElemNode{indices, name}
}

func (node *ArrayElemNode) Eval(scope *Scope) (Variable, error) {
	arr, err := scope.FindVariableDepth(node.name)
	if err != nil {
		return nil, err
	}
	if _, ok := arr.(*VarArray); !ok {
		return nil, fmt.Errorf("%w: '%s'", ErrNotAVectorType, node.name)
	}
	indices, err := convertSizes(node.indices, scope)
	if err != nil {
		return nil, err
	}
	res, _, err := arr.(*VarArray).At(indices)
	return res, err
}
