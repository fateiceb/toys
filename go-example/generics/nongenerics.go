package main

import "fmt"

func SumInts(m map[string]int) int {
	var s int
	for _, v := range m {
		s += v
	}
	return s
}
func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}
func main() {
	mInts := map[string]int{
		"a": 1,
		"b": 2,
		"C": 3,
	}
	mFloat := map[string]float64{
		"a": 1.0,
		"b": 2.0,
		"c": 3.0,
	}
	fmt.Println("IntsSum", SumInts(mInts))
	fmt.Println("floatSum", SumFloats(mFloat))
}
