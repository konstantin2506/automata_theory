package ast

import (
	"reflect"
	"testing"
)

func TestRepeatNode_String(t *testing.T) {
	rnode := RepeatNode{nil, 10}
	wait := "Repeat:10"
	if rnode.String() != wait {
		t.Errorf("incorrect string, wait: %s, got: %s", wait, rnode.String())
	}
}

func TestRepeatNode_Type(t *testing.T) {
	rnode := RepeatNode{nil, 10}
	wait := Repeat
	if rnode.Type() != wait {
		t.Errorf("incorrect string, wait: %d, got: %d", wait, rnode.Type())
	}
}

func TestRepeatNode_Children(t *testing.T) {
	rnode := RepeatNode{&KleeneNode{&CharNode{'c', 0}}, 10}
	res := rnode.Children()
	waitLen := 1
	if len(res) != waitLen {
		t.Errorf("len of children incorrect: %d, wait: %d", len(res), waitLen)
	}
	if reflect.TypeOf(res[0]) != reflect.TypeOf(&KleeneNode{}) {
		t.Errorf("incorrect child type: %T, expect %T", res[0], &KleeneNode{})
	}
}
