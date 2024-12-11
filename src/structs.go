package forum

import (
	"database/sql"
	"time"
)

type Page struct {
	DB *sql.DB
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
	Did         int
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
	Did         int
	Liked       bool
	Created_at  int
	Modified_at int
	Duration    string
	User
}

const TOKEN_AGE = 24 * time.Hour
