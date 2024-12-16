package forum

import (
	"net/http"
	"strings"
	"time"
)

func (data Page) CommentsStore(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	user_id, post_id, body := req.FormValue("user_id"), req.FormValue("post_id"), req.FormValue("body")

	if _, ok := data.Cach[data.Id]; ok {
		errcookie := http.Cookie{
			Name:     "errors",
			Value:    "you are on cooldown! wait for a bit and try again ^_^ .",
			Path:     strings.Split(req.Referer(), "8080")[1],
			MaxAge:   1,
			HttpOnly: true, // Secure the cookie, not accessible by JS
		}
		http.SetCookie(res, &errcookie)
		http.Redirect(res, req, "/posts/"+post_id, http.StatusFound)
		return
	}

	if len(body) < 1 || len(body) > 5000 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	query := "INSERT INTO comments VALUES (NULL,?,?,?,?,?)"
	_, err := data.DB.Exec(query, user_id, post_id, body, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}

	data.Cach[data.Id] = time.Now().Unix()
	http.Redirect(res, req, "/posts/"+post_id, http.StatusFound)
}
