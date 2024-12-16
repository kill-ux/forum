package forum

import (
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

func (data *Page) HandelImage(path string, file multipart.File, fileheader *multipart.FileHeader, res http.ResponseWriter, req *http.Request) string {
	if fileheader.Size > 1000000 {
		return ""
	} else {
		buffer := make([]byte, fileheader.Size)
		_, err := file.Read(buffer)
		if err != nil {
			return ""
		}
		extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"}
		extIndex := slices.IndexFunc(extensions, func(ext string) bool {
			return strings.Contains(fileheader.Filename, ext)
		})
		if extIndex == -1 {
			return ""
		}
		imageName, _ := uuid.NewV4()
		err = os.WriteFile("images/"+path+"/"+imageName.String()+extensions[extIndex], buffer, 0o777)
		if err != nil {
			return ""
		}
		return imageName.String() + extensions[extIndex]
	}
}

func (data Page) PostsStore(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	if _, ok := data.Cach[data.Id]; ok {
		errcookie := http.Cookie{
			Name:     "errors",
			Value:    "you are on cooldown! wait for a bit and try again ^_^ .",
			Path:     "/",
			MaxAge:   1,
			HttpOnly: true, // Secure the cookie, not accessible by JS
		}
		http.SetCookie(res, &errcookie)
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	title, body := req.FormValue("title"), req.FormValue("body")
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if title == "" || len(title) > 400 || body == "" || len(body) > 5000 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	categories := req.Form["category"]
	if len(categories) == 0 {
		categories = []string{"12"}
	}
	file, fileheader, err := req.FormFile("image")
	var image string
	if err == nil {
		image = data.HandelImage("posts", file, fileheader, res, req)
	}
	query := "INSERT INTO posts VALUES (NULL,?,?,?,?,?,?)"
	result, err := data.DB.Exec(query, data.Id, title, body, image, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	} else {
		post_id, _ := result.LastInsertId()
		query = "INSERT INTO posts_categories VALUES (?,?)"
		for _, category := range categories {
			_, err := data.DB.Exec(query, post_id, category)
			if err != nil {
				data.Error(res, http.StatusInternalServerError)
				return
			}
		}
		data.Cach[data.Id] = time.Now().Unix()
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

func (data *Page) PostsInfo(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	strid := req.URL.Path[len("/posts/"):]
	id, atoiErr := strconv.Atoi(strid)
	if atoiErr != nil {
		data.Error(res, http.StatusNotFound)
		return
	}

	//
	page := req.FormValue("page")
	if page == "" {
		page = "1"
	}
	num, err := strconv.Atoi(page)
	if err != nil || num < 1 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	//

	query := "SELECT P.*,U.username,U.image FROM posts P JOIN users U ON P.user_id = U.id WHERE P.id = ?"
	row := data.DB.QueryRow(query, id)
	post := Post{}
	err = row.Scan(&post.Id, &post.User.Id, &post.Title, &post.Body, &post.Image, &post.Created_at, &post.Modified_at, &post.UserName, &post.User.Image)
	if err != nil {
		data.Error(res, http.StatusNotFound)
		return
	}
	post.Duration = CalculDuration(post.Created_at)
	// categories
	query = "SELECT C.* FROM posts_categories PC JOIN categories C ON PC.category_id = C.id Where PC.post_id = ? "
	rows, err := data.DB.Query(query, post.Id)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		category := Categorie{}
		rows.Scan(&category.Id, &category.CatName)
		post.Categories = append(post.Categories, category)
	}
	// likes
	query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 1"
	rowLike := data.DB.QueryRow(query, post.Id)
	rowLike.Scan(&post.Like)

	query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 0"
	rowDisLike := data.DB.QueryRow(query, post.Id)
	rowDisLike.Scan(&post.Dislike)

	query = "SELECT like FROM likes WHERE post_id = ? AND user_id = ?"
	rowDidLike := data.DB.QueryRow(query, post.Id, data.User.Id)
	err = rowDidLike.Scan(&post.Liked)
	if err == nil {
		post.Did = 1
	}
	data.Previous = "/posts/" + strid + "?page=" + strconv.Itoa(num-1)
	if num == 1 {
		data.Previous = "0"
	}
	data.Next = "/posts/" + strid + "?page=" + strconv.Itoa(num+1)
	query = "SELECT C.*,U.username,U.image FROM comments C JOIN users U ON C.user_id = U.id where C.post_id = ? ORDER BY modified_at DESC LIMIT 10 OFFSET " + strconv.Itoa((num-1)*10)
	rows, err = data.DB.Query(query, id)
	for rows.Next() {
		comment := Comment{}
		trash := 0
		rows.Scan(&comment.Id, &comment.User.Id, &trash, &comment.Body, &comment.Created_at, &comment.Modified_at, &comment.UserName, &comment.Image)
		comment.Duration = CalculDuration(comment.Created_at)
		// likes
		query = "SELECT COUNT(*) FROM likes WHERE comment_id = ? AND like = 1"
		rowLike := data.DB.QueryRow(query, comment.Id)
		rowLike.Scan(&comment.Like)

		query = "SELECT COUNT(*) FROM likes WHERE comment_id = ? AND like = 0"
		rowDisLike := data.DB.QueryRow(query, comment.Id)
		rowDisLike.Scan(&comment.Dislike)

		query = "SELECT like FROM likes WHERE comment_id = ? AND user_id = ?"
		rowDidLike := data.DB.QueryRow(query, comment.Id, data.User.Id)
		err := rowDidLike.Scan(&comment.Liked)
		if err == nil {
			comment.Did = 1
		}
		post.Comments = append(post.Comments, comment)
	}
	data.Posts = append(data.Posts, post)
	data.RenderPage("post.html", res)
}

func (data Page) PostsDelete(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	strid := req.URL.Path[len("/posts/delete/"):]
	id, atoiErr := strconv.Atoi(strid)
	if atoiErr != nil {
		data.Error(res, http.StatusNotFound)
		return
	}

	query := "SELECT image FROM posts WHERE id = ? AND user_id = ?"
	row := data.DB.QueryRow(query, id, data.Id)
	image := ""
	err := row.Scan(&image)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	os.Remove("images/posts/" + image)
	query = "DELETE FROM posts WHERE id = ? AND user_id = ?"
	_, err = data.DB.Exec(query, id, data.Id)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/", http.StatusFound)
}

// post update
func (data Page) PostsUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
	}
	post_id, title, body, oldimage := req.FormValue("post_id"), req.FormValue("title"), req.FormValue("body"), req.FormValue("oldimage")
	if post_id == "" || title == "" || len(title) > 400 || body == "" || len(body) > 5000 {
		data.Error(res, http.StatusBadRequest)
		return
	}
	file, fileheader, err := req.FormFile("image")
	var image string
	//

	//
	if err == nil {
		image = data.HandelImage("posts", file, fileheader, res, req)
	}
	//
	query := "SELECT user_id FROM posts WHERE id = ?"
	row := data.DB.QueryRow(query, post_id)
	var user_id int
	err = row.Scan(&user_id)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	if user_id != data.Id {
		data.Error(res, http.StatusBadRequest)
		return
	}
	if image != "" || oldimage == "oldimage" {
		query := "SELECT image FROM posts WHERE id = ? AND user_id = ?"
		row := data.DB.QueryRow(query, post_id, data.Id)
		image2 := ""
		err := row.Scan(&image2)
		if err != nil {
			data.Error(res, http.StatusInternalServerError)
			return
		}
		os.Remove("images/posts/" + image2)
	}
	if image == "" {
		if oldimage == "oldimage" {
			query = "UPDATE posts SET body = ?, title = ?, modified_at = ?, image = ?  WHERE id = ?"
			_, err = data.DB.Exec(query, body, title, time.Now().Unix(), image, post_id)
		} else {
			query = "UPDATE posts SET body = ?, title = ?, modified_at = ?  WHERE id = ?"
			_, err = data.DB.Exec(query, body, title, time.Now().Unix(), post_id)
		}
	} else {
		query = "UPDATE posts SET body = ?, title = ?, image= ?, modified_at = ?  WHERE id = ?"
		_, err = data.DB.Exec(query, body, title, image, time.Now().Unix(), post_id)
	}

	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/posts/"+post_id, http.StatusFound)
}
