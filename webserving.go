package main

import (
	"net"
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"strconv"
	"os"
)

func getIP(req *http.Request, landing string) (string, string, string, string) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)

	//check if splitting host and port returned an error
	if err != nil {
		log.Fatalf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		log.Fatalf("userip: %q is not IP:port", req.RemoteAddr)
		os.Exit(1)
	}

	//this will only be defined when site is accessed via non-anonymous proxy
	//and takes precedence over RemoteAddr
	//header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	log.Printf("IP: %s\n", ip)
	log.Printf("Port: %s\n", port)
	log.Printf("Forwarded For: %s\n", forward)
	log.Printf("Visited: %s\n", landing)
	log.Printf("UA: %s\n", req.UserAgent())

	return ip, port, forward, landing
}

func main() {
	srvPort := strconv.Itoa(3333)

	//instantiate a new router
	r := httprouter.New()

	//add a handler to echo the IP back to the user
	r.GET("/ip", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "<h1>Find Your IP</h1>")
		ip, port, forward, landing := getIP(r, "/ip")

		fmt.Fprintf(w, "<p>IP: %s</p>", ip)
		fmt.Fprintf(w, "<p>Port: %s</p>", port)
		fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)
		fmt.Fprintf(w, "<p>You Are Visiting: %s</p>", landing)
	})

	//add a handler on /
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Simply write some test data for now
		fmt.Fprint(w, "Welcome!\n")
		_, _, _, _ = getIP(r, "/")
	})

	l, err := net.Listen("tcp", "0.0.0.0:" + srvPort)
	if err != nil {
		log.Fatal(err)
	}

	//start the blocking server loop.
	log.Fatal(http.Serve(l, r))
	//mux := http.NewServeMux()
	//mux.Handle("/", http.FileServer(http.Dir("./static")))
	//http.ListenAndServe(":8080", mux)
	//fmt.Println("Connection Handled")

}