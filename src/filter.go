package forum

import (
	"net/http"
	"strconv"
)

func (data *Page) FilterCategories(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	req.ParseForm()
	categories := req.Form["categories"]
	if len(categories) == 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	query := "SELECT DISTINCT P.*,U.username,U.image FROM posts_categories PC JOIN posts P ON PC.post_id = P.id JOIN users U ON P.user_id = U.id WHERE PC.category_id = " + categories[0]
	for _, categorie := range categories[1:] {
		query += " OR PC.category_id = " + categorie
	}
	data.LoadPosts(req, res, query)
}

func (data *Page) LikedPosts(res http.ResponseWriter, req *http.Request) {
	query := "SELECT DISTINCT P.*,U.username,U.image FROM posts P JOIN users U ON P.user_id = U.id JOIN likes L ON P.id = L.post_id WHERE like = 1 AND L.user_id =" + strconv.Itoa(data.Id)
	data.LoadPosts(req, res, query)
}

func (data *Page) CreatedPosts(res http.ResponseWriter, req *http.Request) {
	query := "SELECT DISTINCT P.*,U.username,U.image FROM posts P JOIN users U ON P.user_id = U.id WHERE U.id = " + strconv.Itoa(data.Id)
	data.LoadPosts(req, res, query)
}
