package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	log "log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Content interface{}
	Code string

}
func name(w http.ResponseWriter,r *http.Request) {
	a := map[string]string{}
	a["a"] = "1"
	a["b"] = "2"
	re := Response{Content:a,Code: "success"}
	x, err := json.Marshal(re)
	if err != nil {
		fmt.Sprint(err)
	}
	w.Header().Set("Content-Type","application/json")
	io.Copy(w,bytes.NewReader(x))
}

func file(w http.ResponseWriter,r *http.Request){
	if r.Method == http.MethodGet {
		w.Write([]byte("wrong message"))
	}
	if r.Method == http.MethodPost {
		file,handle,err := r.FormFile("image")
		log.Println(handle.Filename)
		defer  file.Close()
		if err != nil {
			log.Println(err)
		}
		os.Mkdir("uploadimg", 0777)
		saveFile, err := os.OpenFile("uploadimg/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}
		defer saveFile.Close()
		io.Copy(saveFile, file)
	}
}

func main() {
	http.HandleFunc("/name",name)
	http.HandleFunc("/file",file)
	log.Println(http.ListenAndServe(":10001",nil))
	// Echo instance
	//e := echo.New()
	//
	//// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	//
	//// Routes
	//e.GET("/", hello)
	//e.Static("/static", "static")
	//// Start server
	//e.Logger.Fatal(e.Start(":10000"))
}

// Handler
func hello(c echo.Context) error {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println(111)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	body = body[0:1000]
	return c.HTML(http.StatusOK, string(body))
}
