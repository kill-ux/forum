package forum

import (
	"net/http"
	"strings"
)

func (data Page) CommentsLike(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
	comment_id, like, user_id, _ := req.FormValue("comment_id"), req.FormValue("like"), req.FormValue("user_id"), req.FormValue("post_id")

	if like != "0" && like != "1" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	query := "SELECT like FROM likes WHERE user_id = ? AND comment_id = ? AND post_id ISNULL"
	row := data.DB.QueryRow(query, user_id, comment_id)
	liked := false
	err := row.Scan(&liked)
	if err != nil {
		query = "INSERT INTO likes(like,user_id,comment_id) VALUES (?,?,?)"
		data.DB.Exec(query, like, user_id, comment_id)

	} else {
		if (like == "1" && liked) || (like == "0" && !liked) {
			query = "DELETE FROM likes WHERE user_id = ? AND comment_id = ? AND post_id ISNULL"
			data.DB.Exec(query, user_id, comment_id)
		} else {
			query = "UPDATE likes SET like = ? WHERE user_id = ? AND comment_id = ? AND post_id ISNULL"
			data.DB.Exec(query, like, user_id, comment_id)
		}
	}

	http.Redirect(res, req, strings.Split(req.Referer(), ":8080")[1]+"#comment"+comment_id, http.StatusFound)
}

func (data Page) PostsLike(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
	like, post_id := req.FormValue("like"), req.FormValue("post_id")
	user_id := data.Id
	if like != "0" && like != "1" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	query := "SELECT like FROM likes WHERE user_id = ? AND comment_id ISNULL AND post_id = ?"
	row := data.DB.QueryRow(query, user_id, post_id)
	liked := false
	err := row.Scan(&liked)
	if err != nil {
		query = "INSERT INTO likes(like,user_id,post_id) VALUES (?,?,?)"
		data.DB.Exec(query, like, user_id, post_id)

	} else {
		if (like == "1" && liked) || (like == "0" && !liked) {
			query = "DELETE FROM likes WHERE user_id = ? AND comment_id ISNULL AND post_id = ?"
			data.DB.Exec(query, user_id, post_id)
		} else {
			query = "UPDATE likes SET like = ? WHERE user_id = ? AND comment_id ISNULL AND post_id = ?"
			data.DB.Exec(query, like, user_id, post_id)
		}
	}
	http.Redirect(res, req, strings.Split(req.Referer(), ":8080")[1]+"#post"+post_id, http.StatusFound)
}
