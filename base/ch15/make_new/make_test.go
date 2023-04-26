package make_new

import (
	"fmt"
	"testing"
)

/*
TestVariable 会报错

panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x4f503b]
*/

//func TestVariable(t *testing.T) {
//	var i *int
//	*i = 1
//	t.Log(*i)
//}

func TestNew_int(t *testing.T) {
	var i *int
	i = new(int)
	fmt.Println(*i) // int 的零值为 0
	*i = 1
	fmt.Println(*i) // 1
}

func TestMake_slice(t *testing.T) {
	l := make([]int, 0, 10)
	fmt.Println(l == nil)                           // false
	fmt.Printf("len: %d cap: %d\n", len(l), cap(l)) // len: 0 cap: 10
	l = append(l, 1)
	fmt.Println(l) // [1]
}

func TestNew_slice(t *testing.T) {
	l := new([]int)
	fmt.Println(*l == nil)                            // true
	fmt.Printf("len: %d cap: %d\n", len(*l), cap(*l)) // len: 0 cap: 0
	*l = append(*l, 1)
	fmt.Println(*l) // [1]
}

func TestNew_map(t *testing.T) {
	m := new(map[string]int)
	fmt.Println(*m == nil) // true
	(*m)["key"] = 1        // panic: assignment to entry in nil map
}
