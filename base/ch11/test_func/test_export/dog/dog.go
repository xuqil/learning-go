package dog

import (
	"animal"
	"fmt"
)

type Dog struct {
}

func (Dog) Sleep() {
	fmt.Println("狗狗在睡觉")
}

func (Dog) Eat() {
	animal.Eat("狗粮")
}
