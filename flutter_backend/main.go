package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)
type User struct {
	Id []byte `json:"id"`
	Name string`json:"name"`
}
func main() {
	alice := User{[]byte("heelo"),"alice"}
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		log.Println(alice)
		c.JSON(200, alice)
	})
	r.Run() // listen and se
}
