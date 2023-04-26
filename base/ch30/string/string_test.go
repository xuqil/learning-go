package string

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	ch := make(chan string)
	a := "1"
	go func() {
		i := 0
		for {
			if i%2 == 0 {
				a = "1"
			} else {
				a = "22"
			}
			time.Sleep(time.Millisecond * 1) // 阻止编译器优化
			i++
		}
	}()

	go func() {
		for {
			b := a
			if b != "1" && b != "22" {
				ch <- b
			}
		}
	}()

	for i := 0; i < 10; i++ {
		fmt.Println("Got string: ", <-ch)
	}
}

var s string = "abc"

func TestString2(t *testing.T) {
	a := "hello"
	a = "world"
	fmt.Println(a)
}

//func TestModifyString(t *testing.T) {
//	var s string = "abc"
//	s[0] = '0' // Cannot assign to s[0]
//}
