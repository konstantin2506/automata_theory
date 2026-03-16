package kregex

import "fmt"

type Parser struct {
	str             string
	parenPoses      map[int]int
	openParensStack []int
}

func NewParser(s string) Parser {
	return Parser{
		s,
		nil,
		nil,
	}
}

func (p *Parser) ParseParenths() error {
	if p.parenPoses != nil {
		return fmt.Errorf("logic: second use with same string")
	}
	p.parenPoses = make(map[int]int)
	p.openParensStack = make([]int, 0)
	stack := []int{}

	for i, ch := range p.str {
		switch ch {
		case '(':
			stack = append(stack, i)
			p.openParensStack = append(p.openParensStack, i)
		case ')':
			if len(stack) > 0 {
				opened := stack[len(stack)-1]
				p.parenPoses[opened] = i

				stack = stack[:len(stack)-1]
			} else {
				return fmt.Errorf("runtime: redurant ')' symbol in str at %d", i)
			}
		}
	}
	if len(stack) > 0 {
		return fmt.Errorf("runtime: redurant '(' symbol in str at %d", stack[0])
	}
	return nil
}

func (p *Parser) GetNextParenths(counter int) (int, int, error) {
	if counter >= len(p.openParensStack) {
		return 0, 0, fmt.Errorf("logic: out of range (%d)", counter)
	}
	j := len(p.openParensStack) - counter - 1
	return p.openParensStack[j], p.parenPoses[p.openParensStack[j]], nil
}

func (p *Parser) ParseInside(first, last int) {
	fmt.Println(p.str[first+1 : last])
}
