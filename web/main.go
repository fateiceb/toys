package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.Static("/static", "static")
	// Start server
	e.Logger.Fatal(e.Start(":10000"))
}

// Handler
func hello(c echo.Context) error {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println(111)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	body = body[0:1000]
	return c.HTML(http.StatusOK, string(body))
}
