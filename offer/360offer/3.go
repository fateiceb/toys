package main

import "fmt"

func main() {
	// n 个数
	var n, m int
	for {
		a, _ := fmt.Scan(&n)
		if a == 0 {
			break
		}
		nums := make([]int, 0)
		for i := 0; i < n; i++ {
			_, _ = fmt.Scan(&m)
			nums = append(nums, m)
		}
	}
}
