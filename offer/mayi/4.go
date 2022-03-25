package main

import "fmt"

func main() {
	n := 0
	params := []int{1, 2, 3}
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		for j := 0; j < 3; j++ {
			x := 0
			fmt.Scan(&x)
			params[j] = x
		}
		fmt.Printf("%d\n", delnum(params))
	}
}
func delnum(params []int) int {
	return params[0]
}
