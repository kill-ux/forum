package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	forum "forum/src"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	db.SetMaxOpenConns(10)
	forum.DB = db
	// to close db when panic
	defer func() {
		if err := recover(); err != nil {
			db.Close()
			log.Fatal("Error: ", err)
		}
	}()

	// to close db when ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.Close()
		fmt.Println()
		os.Exit(0)
	}()

	creation, err := os.ReadFile("db/creation.sql")
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(string(creation))
	if err != nil {
		panic(err)
	}
	forum.FillCategoriesDB()
	http.HandleFunc("/", forum.Routers)
	// re
	go func() {
		for {
			time.Sleep(time.Second * 10)
			forum.Mux.Lock()
			for key, value := range forum.Cach {
				current := time.Now().Unix() - value
				if current > 10 {
					delete(forum.Cach, key)
				}
			}
			forum.Mux.Unlock()
		}
	}()

	fmt.Println("http://localhost:8080")
	panic(http.ListenAndServe(":8080", nil))
}
