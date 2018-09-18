package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"net/http"
	"os"
)

func getIP(req *http.Request, landing string) (string, string, string, string, string) {
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

	/*this will only be defined when site is accessed via non-anonymous proxy
	/and takes precedence over RemoteAddr
	header.Get is case-insensitive*/
	forward := req.Header.Get("X-Forwarded-For")
	userAgent := req.UserAgent()

	log.Printf("IP: %s\n", ip)
	log.Printf("Port: %s\n", port)
	log.Printf("Forwarded For: %s\n", forward)
	log.Printf("Visited: %s\n", landing)
	log.Printf("UA: %s\n", userAgent)

	return ip, port, forward, landing, userAgent
}

/*func browserLookup(ua string) string {
	var browser string = "UNKNOWN!!!"
	var android map[string]string
	var ios map[string]string
	var desktop map[string]string

	android = make(map[string]string)
	ios = make(map[string]string)
	desktop = make(map[string]string)

	ios["CriOS"] = "Chrome (iPhone)"
	ios["FxiOS"] = "Firefox/Brave (iPhone)"
	ios["Version"] = "Safari (iPhone)[UNCERTAIN]"
	desktop["Firefox"] = "Firefox (Desktop)"
	desktop["Chrome"] = "Chrome (Desktop)"
	android["Firefox"] = "Firefox (Android)"
	android["Chrome"] = "Chrome (Android)"

	select {
		case s.Contains(ua, "Android"):
			for key, value := range android {
				if s.Contains(ua, key) {
					browser = value
				}
			}
		case s.Contains(ua, "iPhone"):
			for key, value := range ios {
				if s.Contains(ua, key) {
					browser = value
				}
			}
		case s.Contains(ua, "Windows NT"):
			for key, value := range desktop {
				if s.Contains(ua, key) {
					browser = value
				}
			}
		case s.Contains(ua, "Linux x86"):
			for key, value := range desktop {
				if s.Contains(ua, key) {
					browser = value
				}
			}
	}

	return browser
}
*/

func main() {
	//instantiate a new router
	r := httprouter.New()

	//add a handler to echo the IP back to the user
	r.GET("/ip", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "<h1>Find Your IP</h1>")
		ip, port, forward, landing, userAgent := getIP(r, "/ip")

		fmt.Fprintf(w, "<p>IP: %s</p>", ip)
		fmt.Fprintf(w, "<p>Port: %s</p>", port)
		fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)
		fmt.Fprintf(w, "<p>You Are Visiting: %s</p>", landing)
		fmt.Fprintf(w, "<p>Your User Agent is: %s</p>", userAgent)
		fmt.Fprintf(w, "<p>This means you are using the %s browser</p>", "Not Sure Yet")
	})

	//add a handler on /
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var ip string
		ip, _, _, _, _ = getIP(r, "/")
		fmt.Fprintf(w, "<h1>Welcome %s!</h1>Please to be enjoying your stayings!\n", ip)
	})

	//start the blocking server loop.
	log.Fatal(http.Serve(autocert.NewListener("init.tools"), r))
}
