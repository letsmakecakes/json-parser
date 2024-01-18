package stack

import "errors"

var stack []int

func Push(val int) {
	stack = append(stack, val)
}

func Pop() (int, error) {
	if len(stack) == 0 {
		return -1, errors.New("stack is empty")
	}
	tmp := stack[len(stack)-1]
	stack = stack[0 : len(stack)-1]
	return tmp, nil
}

func GetStack() []int {
	return stack
}

func Peek() int {
	return stack[len(stack)-1]
}

func IsEmpty() bool {
	return len(stack) != 0
}
