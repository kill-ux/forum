package forum

import (
	"net/http"
	"time"
)

func (data Page) CommentsStore(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	user_id, post_id, body := req.FormValue("user_id"), req.FormValue("post_id"), req.FormValue("body")
	if len(body) > 5000 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	query := "INSERT INTO comments VALUES (NULL,?,?,?,?,?)"
	_, err := data.DB.Exec(query, user_id, post_id, body, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/posts/"+post_id, http.StatusFound)
}
