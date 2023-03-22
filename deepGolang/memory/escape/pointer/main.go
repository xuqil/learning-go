package main

type Person struct {
	name string
}

func NewPerson(name string) *Person {
	p := &Person{ // 局部变量 p 逃逸到堆
		name: name,
	}
	return p
}

func main() {
	NewPerson("tom")
}
