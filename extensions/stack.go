package extensions

import "fmt"

type Stack struct {
	items []uint16
}

func (s *Stack) Push(data uint16) {
	s.items = append(s.items, data)
}

func (s *Stack) Pop() uint16 {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	var res = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return res
}

func (s *Stack) Top() (uint16, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Count() int {
	return len(s.items)
}

func (s *Stack) Print() {
	for _, item := range s.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
