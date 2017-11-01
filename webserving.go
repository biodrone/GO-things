package main

import ("net/http"
		"database/sql"
		_"github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", mux)
}

/*
// open database
	database, _ := sql.Open("sqlite3", "./TestDatabase.db")

//create table - also creates database if it doesn't already exist
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
	statement.Exec()

//insert
	statement, _ = database.Prepare("INSERT INTO Users (firstname, lastname) VALUES (?, ?)")
	statement.Exec("admin", "Memes")

//update
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
	}*/