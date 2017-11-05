package main

import (
	"net"
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"log"
	"os"
	"golang.org/x/crypto/acme/autocert"
	s "strings"
	"time"
)

type IPLog struct {
	ip, port, forward, landing, userAgent string
	createdDate time.Time
}

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

func browserLookup(ua string) string {
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

func insertLogRecord(ip string, port string, forward string, landing string, userAgent string) {
	database, err1:= sql.Open("sqlite3", "./LogDatabase.db")
	if err1 != nil {
		log.Fatal(err1)
	}

	statement, err2 := database.Prepare("INSERT INTO Logs (ip, port, forward, landing, user_agent) VALUES (?, ?, ?, ?, ?)")
	if err2 != nil {
		log.Fatal(err2)
	}
	statement.Exec(ip, port, forward, landing, userAgent)
	database.Close()
}

func getLogRecords() []IPLog{
	database, _ := sql.Open("sqlite3", "./LogDatabase.db")

	rows, _ := database.Query("SELECT * FROM Logs")
	var logs []IPLog
	for rows.Next() {
		var ipLog IPLog
		rows.Scan(&ipLog.ip, &ipLog.port, &ipLog.forward, &ipLog.landing, &ipLog.userAgent, &ipLog.createdDate)
		logs = append(logs, ipLog)
	}
	database.Close()
	return logs
}

func getLogsCount() int {
	database, _ := sql.Open("sqlite3", "./LogDatabase.db")

	var numRows int
	countRows, _ := database.Query("SELECT COUNT(*) FROM Logs")
	countRows.Scan(&numRows)
	return numRows
}

/*
// open database
	database, _ := sql.Open("sqlite3", "./TestDatabase.db")

//create table - also creates database if it doesn't already exist
	statement, _ := database.Prepare("CREATE TABLE Logs (ip TEXT, port TEXT, forward TEXT, landing TEXT, user_agent TEXT, created_date DATETIME DEFAULT(CURRENT_TIMESTAMP)")
	statement.Exec()

//insert
	statement, _ := database.Prepare("INSERT INTO Users (firstname, lastname) VALUES (?, ?)")
	statement.Exec("admin", "Memes")

//update
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
	}*/

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
		fmt.Fprintf(w, "<p>This means you are using the %s browser</p>", browserLookup(userAgent))
	})

	//add a handler on /
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var ip string
		ip, _, _, _, _= getIP(r, "/")
		fmt.Fprintf(w, "<h1>Welcome %s!</h1>Please to be enjoying your stayings!\n", ip)
	})

	//start the blocking server loop.
	log.Fatal(http.Serve(autocert.NewListener("jjgo.init.tools"), r))
	//mux := http.NewServeMux()
	//mux.Handle("/", http.FileServer(http.Dir("./static")))
	//http.ListenAndServe(":8080", mux)
	//fmt.Println("Connection Handled")
}
