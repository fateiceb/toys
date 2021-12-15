package main
import "fmt"
func SumIntsOrFloats[K comparable,V int | float64](m map[K]V) V {
	var s V
	for _,v := range m {
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
	fmt.Println("sum",SumIntsOrFloats[string,int](mInts)) // []中声明函数中应该被替换的参数类型
	fmt.Println("IntsSum", SumIntsOrFloats(mInts))
	fmt.Println("floatSum", SumIntsOrFloats(mFloat))
}
