package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	forum "forum/src"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	creation, err := os.ReadFile("db/creation.sql")
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Exec(string(creation))
	data := forum.Page{DB: db, Cach: map[int]int64{}}
	data.FillCategories()
	http.HandleFunc("/", data.Routers)
	http.HandleFunc("/css/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			data.Error(res, http.StatusMethodNotAllowed)
			return
		}
		_, err := os.ReadFile(req.URL.Path[1:])
		if err != nil {
			data.Error(res, http.StatusNotFound)
			return
		}
		http.StripPrefix("/css/", http.FileServer(http.Dir("css"))).ServeHTTP(res, req)
	})
	http.HandleFunc("/images/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			data.Error(res, http.StatusMethodNotAllowed)
			return
		}
		_, err := os.ReadFile(req.URL.Path[1:])
		if err != nil {
			data.Error(res, http.StatusInternalServerError)
			return
		}
		http.FileServer(http.Dir(".")).ServeHTTP(res, req)
	})

	// re
	go func() {
		for {
			time.Sleep(time.Second * 10)
			for key, value := range data.Cach {
				current := time.Now().Unix() - value
				if current > 10 {
					delete(data.Cach, key)
				}
			}
		}
	}()

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
