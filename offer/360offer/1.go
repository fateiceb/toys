package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//输入
	var cnt int
	reader := bufio.NewReader(os.Stdin)
	srcdata, _ := reader.ReadString('\n')
	targetdata, _ := reader.ReadString('\n')
	//转换为字节串
	src := []byte(srcdata)
	target := []byte(targetdata)
	//统计次数
	ms := make(map[byte]int, 0)
	mt := make(map[byte]int, 0)
	for i := 0; i < len(src); i++ {
		ms[src[i]]++
		mt[target[i]]++
	}
	//目标串和原始串字符数量不一致，只能删除

	//匹配
	for i := range src {
		if src[i] != target[i] {
			cnt++
		}
	}
	cnt = cnt / 2
	if (abs(ms['A']-mt['A']))%2 != 0 {
		cnt++
	}
	fmt.Println(cnt)
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/*
func ts0(src, dst string) int {
    res := 0
    n := 0
    srcbyte := []byte(src)
    dstbyte := []byte(dst)
    srcmap := make(map[byte]int)
    dstmap := make(map[byte]int)
    for i := 0; i < len(src); i++ {
        srcmap[src[i]]++
        dstmap[dst[i]]++
    }
    for i := 0; i < len(src) && srcmap['A'] < dstmap['A']; i++ {
        if src[i] == 'T' && dst[i] == 'A' {
            srcbyte[i] = byte('A')
            res++
            srcmap['A']++
            srcmap['T']--
        }
    }
    for i := 0; i < len(src) && srcmap['A'] > dstmap['A']; i++ {
        if src[i] == 'A' && dst[i] == 'T' {
            srcbyte[i] = byte('T')
            res++
            srcmap['A']--
            srcmap['T']++
        }
    }
    for i := 0; i < len(srcbyte); i++ {
        if srcbyte[i] != dstbyte[i] {
            n++
        }
    }
    return res + n/2
}

*/
