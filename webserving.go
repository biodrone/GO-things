package main

import (
	"net/http")

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", mux)
}
