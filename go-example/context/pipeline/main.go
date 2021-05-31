package main

import (
	"fmt"
	"sync"
)

func main() {
	c := gen(2,3)
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func(){
		for _,n := range nums {
			out <- n
			fmt.Println(n)
		}
		close(out)
	}()
	return out
}
func sq(in <-chan int) <-chan int{
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int{
	var wg sync.WaitGroup
	out := make(chan int)
	output := func(c <-chan int) {
		for n := range c {
			out <-n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _,c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}



func mergecopy(cs ...<-chan int) <- chan int{
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	output := func(n <- chan int) {
		for i := range n {
			out <- i
		}
		wg.Done()
	}
	for _,i := range cs{
		go output(i)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}