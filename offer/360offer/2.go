package main

import (
	"fmt"
	"io"
)

func main() {
	var n, q int
	fmt.Scanf("%d%d", &n, &q)
	arr := make([]int, n+1)

	for {
		var l, r, cnt int
		_, err := fmt.Scanf("%d%d", &l, &r)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		//翻转
		for i := l; i < r+1; i++ {
			if arr[i] == 1 {
				arr[i] = 0
			} else {
				arr[i] = 1
			}
		}
		//统计
		for i := 1; i <= n; i++ {
			if arr[i] == 0 {
				cnt++
			}
		}
		fmt.Println(cnt)
	}
}
