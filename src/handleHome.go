package forum

import (
	"net/http"
	"strconv"
	"time"
)

func CalculDuration(Created_at int) string {
	dur := int(time.Now().Unix() - int64(Created_at))
	if dur < 60 {
		// s
		return strconv.Itoa(dur) + " seconds"
	} else if dur < 3600 {
		// m
		return strconv.Itoa(dur/60) + " minutes"
	} else if dur < 86400 {
		// H
		return strconv.Itoa(dur/3600) + " hours"
	} else if dur < 2592000 {
		// D
		return strconv.Itoa(dur/86400) + " days"
	} else if dur < 31104000 {
		// M
		return strconv.Itoa(dur/2592000) + " months"
	} else {
		// y
		return strconv.Itoa(dur/31104000) + " years"
	}
}

func (data *Page) LoadPosts(req *http.Request, res http.ResponseWriter, query string) {
	page := req.FormValue("page")
	if page == "" {
		page = "1"
	}
	num, err := strconv.Atoi(page)
	if err != nil || num < 1 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	url := req.URL.Query()
	url.Set("page", strconv.Itoa(num-1))
	req.URL.RawQuery = url.Encode()
	data.Previous = req.URL.Path + "?" + req.URL.RawQuery
	if num == 1 {
		data.Previous = "0"
	}
	url.Set("page", strconv.Itoa(num+1))
	req.URL.RawQuery = url.Encode()
	data.Next = req.URL.Path + "?" + req.URL.RawQuery
	req.URL.RawQuery = url.Encode()
	tempQuery := "SELECT count(*) nb FROM posts"
	nb := 0
	row := data.DB.QueryRow(tempQuery)
	err = row.Scan(&nb)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	if nb%10 != 0 {
		nb += 10
	}
	if num != nb/10 {
		data.Current = true
	} else {
		data.Current = false
	}
	url.Set("page", strconv.Itoa(nb/10))
	req.URL.RawQuery = url.Encode()
	data.Last = req.URL.Path + "?" + req.URL.RawQuery
	query += " ORDER BY modified_at DESC LIMIT 10 OFFSET " + strconv.Itoa((num-1)*10)
	rows, err := data.DB.Query(query)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		post := Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Title, &post.Body, &post.Image, &post.Created_at, &post.Modified_at, &post.UserName, &post.User.Image)

		post.Duration = CalculDuration(post.Created_at)
		// likes
		query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 1"
		rowLike := data.DB.QueryRow(query, post.Id)
		rowLike.Scan(&post.Like)

		query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 0"
		rowDisLike := data.DB.QueryRow(query, post.Id)
		rowDisLike.Scan(&post.Dislike)

		query = "SELECT like FROM likes WHERE post_id = ? AND user_id = ?"
		rowDidLike := data.DB.QueryRow(query, post.Id, data.User.Id)
		err := rowDidLike.Scan(&post.Liked)
		if err == nil {
			post.Did = 1
		}
		query = "SELECT C.* FROM posts_categories PC JOIN categories C ON PC.category_id = C.id Where PC.post_id = ? "
		rows, err := data.DB.Query(query, post.Id)
		if err != nil {
			data.Error(res, http.StatusBadRequest)
			return
		}
		for rows.Next() {
			category := Categorie{}
			rows.Scan(&category.Id, &category.CatName)
			post.Categories = append(post.Categories, category)
		}
		//
		trash := 0
		trash1 := 0
		query = `	SELECT 
						C.*, 
						U.username, 
						U.image, 
						L.nb_likes 
					FROM 
						(SELECT comment_id, COUNT(*) AS nb_likes 
						FROM likes 
						GROUP BY comment_id 
						HAVING comment_id IS NOT NULL) AS L
					JOIN comments C ON L.comment_id = C.id 
						AND C.post_id = ?
					JOIN users U ON C.user_id = U.id
					ORDER BY L.nb_likes DESC LIMIT 1
		`

		row := data.DB.QueryRow(query, post.Id)
		row.Scan(&post.TopComment.Id, &post.TopComment.User.Id, &trash, &post.TopComment.Body, &post.TopComment.Created_at, &post.TopComment.Modified_at, &post.TopComment.UserName, &post.TopComment.Image, &trash1)
		//
		data.Posts = append(data.Posts, post)
	}
	data.RenderPage("index.html", res)
}

func (data *Page) HandleHome(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	query := "SELECT P.*,U.username,U.image FROM posts P JOIN users U ON P.user_id = U.id "
	data.LoadPosts(req, res, query)
}
