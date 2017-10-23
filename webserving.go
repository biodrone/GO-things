package main

import (
	"fmt"
	"net/http")

func index_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Dis is da good shit!")
	fmt.Println("Connection to main page!")
}

func about_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About da good shit!")
	fmt.Println("Connection to about page!")
}

func main() {
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/about", about_handler)
	http.ListenAndServe(":8080", nil)
}
