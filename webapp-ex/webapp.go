package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(handleRequest))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "<h1>Hello World!</h1>")
}
