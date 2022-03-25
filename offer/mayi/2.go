// 本题为考试单行多行输入输出规范示例，无需提交，不计分。
package main

import (
	"fmt"
	"unicode"
)

func main() {
	n := 0
	fmt.Scan(&n)
	m := make(map[string]bool)
	var s string
	for {
		n, _ := fmt.Scan(&s)
		if n == 0 {
			break
		} else {

			fmt.Printf("%d\n%s", n, s)
		}
	}
}
func legal(s string, m map[string]bool) string {
	success := "registration complete"
	err1 := "illegal length"
	err2 := "illegal charactor"
	err3 := "acount existed"
	if len(s) > 12 || len(s) < 6 {
		return err1
	}
	for _, c := range s {
		if !unicode.IsLetter(c) {
			return err2
		}
	}
	if m[s] {
		return err3
	} else {
		m[s] = true
	}
	return success
}
