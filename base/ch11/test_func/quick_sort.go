package quick_sort

// QuickSort 快速排序， arr 为 int 类型的切片
func QuickSort(arr []int) {
	quickSortR(arr, 0, len(arr)-1)
}

func quickSortR(arr []int, p, r int) {
	if p >= r {
		return
	}
	q := partition(arr, p, r)
	quickSortR(arr, p, q-1)
	quickSortR(arr, q+1, r)
}

func partition(arr []int, p, r int) int {
	i, j := p, r-1
	for i <= j {
		if arr[i] < arr[r] {
			i++
			continue
		}
		if arr[j] >= arr[r] {
			j--
			continue
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[i], arr[r] = arr[r], arr[i]
	return i
}
