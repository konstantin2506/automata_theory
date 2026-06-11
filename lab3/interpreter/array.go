package interpreter

import (
	"errors"
	"fmt"
)

var (
	ErrArraySizes      = errors.New("array sizes are incorrect")
	ErrArrayIndices    = errors.New("array indices count differ")
	ErrArrayOutOfRange = errors.New("array index is out of range")
)

type VarArray[T any] struct {
	sizes []int
	data  []T
	dataT VarT
}

func (v *VarArray[T]) Assign(indices []int, value T) error {
	if len(indices) != len(v.sizes) {
		return fmt.Errorf("%w: sizes: %d, indices: %d", ErrArrayIndices, len(v.sizes), len(indices))
	}

	resultIndex := indices[len(indices)-1]
	if resultIndex > v.sizes[len(v.sizes)-1] {
		return fmt.Errorf("%w: idx[%d] = %d, dim[%d] = %d", ErrArrayOutOfRange, len(indices)-1, resultIndex, len(indices)-1, v.sizes[len(indices)-1])
	}
	for i := len(indices) - 1; i >= 0; i-- {
		if v.sizes[i] < indices[i] {
			return fmt.Errorf("%w: idx[%d] = %d, dim[%d] = %d", ErrArrayOutOfRange, i, indices[i], i, v.sizes[i])
		}
		resultIndex = resultIndex*v.sizes[i] + indices[i]
	}
	v.data[resultIndex] = value
	return nil
}

func (v *VarArray[T]) Type() VarT {
	return v.dataT
}

func newVariableArray[T any](sizes []int, value T, dataT VarT) (*VarArray[T], error) {
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
	data := make([]T, dataSize)
	for i := range dataSize {
		data[i] = value
	}
	v := &VarArray[T]{sizesCopy, data, dataT}
	return v, nil
}

func NewBoolArray(sizes []int, value bool) (Variable, error) {
	arr, err := newVariableArray(sizes, value, Bool)
	if err != nil {
		return nil, err
	}
	return arr, err
}

func NewIntArray(sizes []int, value int) (Variable, error) {
	arr, err := newVariableArray(sizes, value, Int)
	if err != nil {
		return nil, err
	}
	return arr, err
}
