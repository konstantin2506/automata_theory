package kregex

import (
	"reflect"
	"testing"
)

func TestOrNode_String(t *testing.T) {
	rnode := OrNode{nil, false}
	wait := "Or"
	if rnode.String() != wait {
		t.Errorf("incorrect string, wait: %s, got: %s", wait, rnode.String())
	}
}

func TestOrNode_Type(t *testing.T) {
	rnode := OrNode{nil, false}
	wait := Or
	if rnode.Type() != wait {
		t.Errorf("incorrect string, wait: %d, got: %d", wait, rnode.Type())
	}
}

func TestOrNode_Children(t *testing.T) {
	childs := []Node{}
	childs = append(childs, &KleeneNode{})
	childs = append(childs, &CharNode{})
	ornode := OrNode{childs, false}
	res := ornode.Children()
	waitLen := 2
	if len(res) != waitLen {
		t.Errorf("len of children incorrect: %d, wait: %d", len(res), waitLen)
	}
	if reflect.TypeOf(res[0]) != reflect.TypeOf(&KleeneNode{}) {
		t.Errorf("incorrect child[0] type: %T, expect %T", res[0], &KleeneNode{})
	}
	if reflect.TypeOf(res[1]) != reflect.TypeOf(&CharNode{}) {
		t.Errorf("incorrect child[0] type: %T, expect %T", res[0], &CharNode{})
	}
}
