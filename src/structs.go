package forum

import (
	"database/sql"
	"sync"
	"time"
)

type Page struct {
	User
	Posts      []Post
	Categories []Categorie
	PageError  struct {
		Code int
		Msg  string
	}
	Previous string
	Next     string
	Last     string
	Current  bool
}

var (
	DB   *sql.DB
	Cach = map[int]int64{}
	Mux sync.Mutex
)

type Categorie struct {
	Id      int
	CatName string
}

type User struct {
	Id         int
	UserName   string
	Email      string
	Password   string
	Role       string
	Token      any
	Token_Exp  int
	Image      string
	Created_at string
	Log        int
	ErrorEmail string
}

type Post struct {
	Id          int
	Image       string
	Comments    []Comment
	TopComment  Comment
	Title       string
	Body        string
	Like        int
	Dislike     int
	Did         bool
	Liked       bool
	Created_at  int
	Modified_at int
	Duration    string
	User
	Categories []Categorie
}

type Comment struct {
	Id          int
	Body        string
	Like        int
	Dislike     int
	Did         bool
	Liked       bool
	Created_at  int
	Modified_at int
	Duration    string
	User
}

const TOKEN_AGE = 24 * time.Hour
