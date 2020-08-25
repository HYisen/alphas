package utility

type Stack struct {
	data []interface{}
}

func (s *Stack) Top() interface{} {
	//fmt.Printf("top stack size=%d\n", len(s.data))
	return s.data[len(s.data)-1]
}

// Want get the removed item? Use Top before Pop.
func (s *Stack) Pop() {
	s.data = s.data[:len(s.data)-1]
	//fmt.Printf("popped stack size=%d\n", len(s.data))
}

func (s *Stack) Push(neo interface{}) {
	//fmt.Printf("pushing stack size=%d\n", len(s.data))
	//fmt.Println(s.data)
	s.data = append(s.data, neo)
	//fmt.Printf("pushed stack size=%d\n", len(s.data))
}
