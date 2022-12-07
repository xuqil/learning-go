package animal

import "fmt"

type Animal interface {
	Sleep()
	Eat()
}

func Eat(food string) {
	fmt.Printf("吃%s", food)
}
