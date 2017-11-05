package main

import ("net/http"
		"database/sql"
		_"github.com/mattn/go-sqlite3"
		"time"
		"log"
)

type IPLog struct {
	ip, port, forward, landing, userAgent string
	createdDate time.Time
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", mux)
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