package cat

import (
	"animal"
	"fmt"
)

type Cat struct {
}

func (Cat) Sleep() {
	fmt.Println("猫咪在睡觉")
}

func (Cat) Eat() {
	animal.Eat("猫粮")
}
