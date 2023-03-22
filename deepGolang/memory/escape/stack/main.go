package main

func Slice8192() {
	_ = make([]int, 8192) // = 64KB
}

func Slice8193() {
	_ = make([]int, 8193) // > 64KB
}

func SliceUnknown(n int) {
	_ = make([]int, n) // 不确定大小
}

func main() {
	Slice8192()
	Slice8193()
	SliceUnknown(1)
}
