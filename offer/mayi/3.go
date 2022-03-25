package main

import (
	"fmt"
	"sort"
)

func main() {
	n := 0
	params := []int{1, 2, 3, 4, 5}
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		for j := 0; j < 5; j++ {
			x := 0
			fmt.Scan(&x)
			params[j] = x
		}
		fmt.Printf("%d\n", delnum(params))
	}
}
func delnum(arr []int) int {
	sort.Ints(arr)
	arr[i] = -1

}
