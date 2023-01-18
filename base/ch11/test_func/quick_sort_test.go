package quick_sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestQuickSort(t *testing.T) {
	testCases := []struct {
		name string
		arr  []int

		// 测试预期的结果
		want []int
	}{
		{
			name: "normal",
			arr:  []int{4, 1, 2, 3, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "duplicate",
			arr:  []int{6, 6, 2, 3, 5},
			want: []int{2, 3, 5, 6, 6},
		},
		{
			name: "empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			QuickSort(tc.arr)
			if !equal(tc.arr, tc.want) {
				t.Errorf(`want: %v actual: %v`, tc.want, tc.arr)
			}
		})
	}
}

// 带重复数字的切片
func TestQuickSortDuplicate(t *testing.T) {
	numbers := []int{6, 6, 2, 3, 5}
	QuickSort(numbers)
	want := []int{2, 3, 5, 6, 6}
	if !equal(numbers, want) {
		t.Errorf(`want: %v actual: %v`, want, numbers)
	}
}

// 空切片
func TestQuickSortEmpty(t *testing.T) {
	var numbers []int
	QuickSort(numbers)
	var want []int
	if !equal(numbers, want) {
		t.Errorf(`want: %v actual: %v`, want, numbers)
	}
}
func equal(arr1, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}

func TestQuickSort_ByRandom(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		arr, want := randomSlice(rng)
		QuickSort(arr)
		if !equal(arr, want) {
			t.Errorf(`want: %v actual: %v`, want, arr)
		}
	}
}

func randomSlice(rng *rand.Rand) ([]int, []int) {
	n := rng.Intn(25)
	var numbers []int
	for i := 0; i < n; i++ {
		numbers = append(numbers, rng.Intn(100))
	}

	sortedNumbers := make([]int, len(numbers))
	copy(sortedNumbers, numbers)
	// 这里使用标准库的排序得到正确的排序结果
	sort.Ints(sortedNumbers)
	return numbers, sortedNumbers
}

func TestCoverage(t *testing.T) {
	testCases := []struct {
		name string
		arr  []int

		// 测试预期的结果
		want []int
	}{
		{
			name: "normal",
			arr:  []int{4, 1, 2, 3, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "duplicate",
			arr:  []int{6, 6, 2, 3, 5},
			want: []int{2, 3, 5, 6, 6},
		},
		{
			name: "empty",
		},
		{
			name: "reverse",
			arr:  []int{5, 4, 3, 2, 1},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			QuickSort(tc.arr)
			if !equal(tc.arr, tc.want) {
				t.Errorf(`want: %v actual: %v`, tc.want, tc.arr)
			}
		})
	}
}

func BenchmarkQuickSort(b *testing.B) {
	arr := []int{4, 1, 2, 3, 5}
	for i := 0; i < b.N; i++ {
		QuickSort(arr)
	}
}

func BenchmarkQuickSort_Parallel(b *testing.B) {
	arr := []int{4, 1, 2, 3, 5}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起执行，循环一共执行 b.N 次
			QuickSort(arr)
		}
	})
}
