package main

import (
	"fmt"
	"gee"
	"net/http"
	"time"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	r.GET("/gettime", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Time is : %q\n", time.Now().Format("2006-01-03 15:04:05"))
	})

	r.Run(":9998")
}
