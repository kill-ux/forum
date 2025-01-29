package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

const maxFileSize = 1000000 // 1MB file size limit

// Helper function to handle image uploads
func (data *Page) HandleImage(path string, file multipart.File, fileheader *multipart.FileHeader) string {
	if fileheader.Size > maxFileSize {
		return ""
	}

	buffer := make([]byte, fileheader.Size)
	_, err := file.Read(buffer)
	if err != nil {
		return ""
	}

	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"}
	extIndex := slices.IndexFunc(extensions, func(ext string) bool {
		return strings.HasSuffix(fileheader.Filename, ext)
	})
	if extIndex == -1 {
		return ""
	}

	imageName, _ := uuid.NewV4()
	fmt.Println(os.ErrExist)
	err = os.WriteFile("images/"+path+"/"+imageName.String()+extensions[extIndex], buffer, 0o644) // Safer permissions
	if err != nil {
		return ""
	}
	return imageName.String() + extensions[extIndex]
}

// // Helper function to redirect to a specific path after performing an action
// func redirectWithReferer(res http.ResponseWriter, req *http.Request, postID string) {
// 	redirectURL := getPathFromReferer(req) + "#post" + postID
// 	http.Redirect(res, req, redirectURL, http.StatusFound)
// }

// Function to store a new post
func (data *Page) PostsStore(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	// Check for cooldown
	Mux.Lock()
	_, onCooldown := Cach[data.Id]
	Mux.Unlock()
	if onCooldown {
		http.SetCookie(res, &http.Cookie{
			Name:     "errors",
			Value:    "you are on cooldown! wait for a bit and try again ^_^ .",
			Path:     "/",
			MaxAge:   1,
			HttpOnly: true,
		})
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Validate post inputs
	title, body := strings.TrimSpace(req.FormValue("title")), strings.TrimSpace(req.FormValue("body"))
	if title == "" || len(title) > 400 || body == "" || len(body) > 5000 {
		data.Error(res, http.StatusBadRequest)
		return
	}

	// Handle categories
	categories := req.Form["category"]
	if len(categories) == 0 {
		categories = []string{"12"}
	}

	// Process image if exists
	file, fileheader, err := req.FormFile("image")
	var image string
	if err == nil {
		image = data.HandleImage("posts", file, fileheader)
	}

	// Insert post into database
	query := "INSERT INTO posts VALUES (NULL,?,?,?,?,?,?)"
	result, err := DB.Exec(query, data.Id, title, body, image, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}

	// Insert post categories
	postID, _ := result.LastInsertId()
	query = "INSERT INTO posts_categories VALUES (?,?)"
	for _, category := range categories {
		_, err := DB.Exec(query, postID, category)
		if err != nil {
			data.Error(res, http.StatusInternalServerError)
			return
		}
	}

	// Update cooldown cache
	Mux.Lock()
	Cach[data.Id] = time.Now().Unix()
	Mux.Unlock()

	http.Redirect(res, req, "/", http.StatusFound)
}

func (data *Page) CommentsInfo(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	strID := req.URL.Path[len("/posts/"):strings.Index(req.URL.Path,"/comments")]
	postID, err := strconv.Atoi(strID)
	if err != nil {
		data.Error(res, http.StatusNotFound)
		return
	}

	// Handle pagination for comments
	page := req.FormValue("page")
	if page == "" {
		page = "1"
	}
	
	num, err := strconv.Atoi(page)
	if err != nil || num < 1 {
		data.Error(res, http.StatusBadRequest)
		return
	}

	// Pagination for comments
	data.Previous = "/posts/" + strID + "?page=" + strconv.Itoa(num-1)
	if num == 1 {
		data.Previous = "0"
	}
	data.Next = "/posts/" + strID + "?page=" + strconv.Itoa(num+1)


	// Fetch comments for the post
	query := `SELECT C.*, U.username, U.image,
		(SELECT COUNT(*) FROM likes WHERE comment_id = C.id AND like = 1) as comment_likes,
		(SELECT COUNT(*) FROM likes WHERE comment_id = C.id AND like = 0) as comment_dislikes,
		COALESCE((SELECT like FROM likes WHERE comment_id = C.id AND user_id = ?), "") as did
	FROM comments C JOIN users U ON C.user_id = U.id WHERE C.post_id = ?
	ORDER BY modified_at DESC LIMIT 10 OFFSET ` + strconv.Itoa((num-1)*10)
	rows, err := DB.Query(query, data.Id, postID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		data.Error(res, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Process comments
	var Response struct {
		Comments []Comment
		Previous string
		Next string
	}

	for rows.Next() {
		var comment Comment
		var trash int
		var did string
		err = rows.Scan(&comment.Id, &comment.User.Id, &trash, &comment.Body, &comment.Created_at, &comment.Modified_at, &comment.UserName, &comment.Image, &comment.Like, &comment.Dislike, &did)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
		}
		comment.Duration = CalculDuration(comment.Created_at)
		if did != "" {
			comment.Did = true
			comment.Liked = (did == "1")
		}
		Response.Comments = append(Response.Comments, comment)
	}
	Response.Previous = data.Previous
	Response.Next = data.Next
	json.NewEncoder(res).Encode(Response)
}

// Function to show the details of a post
func (data *Page) PostsInfo(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	strID := req.URL.Path[len("/posts/"):]
	postID, err := strconv.Atoi(strID)
	if err != nil {
		data.Error(res, http.StatusNotFound)
		return
	}

	// Handle pagination for comments
	page := req.FormValue("page")
	if page == "" {
		page = "1"
	}
	num, err := strconv.Atoi(page)
	if err != nil || num < 1 {
		data.Error(res, http.StatusBadRequest)
		return
	}

	// Fetch post details from database
	query := "SELECT P.*, U.username, U.image FROM posts P JOIN users U ON P.user_id = U.id WHERE P.id = ?"
	row := DB.QueryRow(query, postID)
	var post Post
	err = row.Scan(&post.Id, &post.User.Id, &post.Title, &post.Body, &post.Image, &post.Created_at, &post.Modified_at, &post.UserName, &post.User.Image)
	if err != nil {
		data.Error(res, http.StatusNotFound)
		return
	}
	post.Duration = CalculDuration(post.Created_at)

	// Fetch categories for the post
	query = "SELECT C.* FROM posts_categories PC JOIN categories C ON PC.category_id = C.id WHERE PC.post_id = ?"
	rows, err := DB.Query(query, post.Id)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var category Categorie
		err = rows.Scan(&category.Id, &category.CatName)
		if err != nil {
			fmt.Println(err)
		}
		post.Categories = append(post.Categories, category)
	}
	rows.Close()

	// Fetch likes and dislikes for the post
	query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 1"
	rowLike := DB.QueryRow(query, post.Id)
	if err = rowLike.Scan(&post.Like); err != nil {
		fmt.Println(err)
	}

	query = "SELECT COUNT(*) FROM likes WHERE post_id = ? AND like = 0"
	rowDisLike := DB.QueryRow(query, post.Id)
	if err = rowDisLike.Scan(&post.Dislike); err != nil {
		fmt.Println(err)
	}

	// Check if the user has liked or disliked the post
	query = "SELECT like FROM likes WHERE post_id = ? AND user_id = ?"
	rowDidLike := DB.QueryRow(query, post.Id, data.User.Id)
	err = rowDidLike.Scan(&post.Liked)
	if err == nil {
		post.Did = true
	} else if err != sql.ErrNoRows {
		log.Println(err)
	}

	// Pagination for comments
	data.Previous = "/posts/" + strID + "?page=" + strconv.Itoa(num-1)
	if num == 1 {
		data.Previous = "0"
	}
	data.Next = "/posts/" + strID + "?page=" + strconv.Itoa(num+1)

	// Fetch comments for the post
	query = `SELECT C.*, U.username, U.image,
		(SELECT COUNT(*) FROM likes WHERE comment_id = C.id AND like = 1) as comment_likes,
		(SELECT COUNT(*) FROM likes WHERE comment_id = C.id AND like = 0) as comment_dislikes,
		COALESCE((SELECT like FROM likes WHERE comment_id = C.id AND user_id = ?), "") as did
	FROM comments C JOIN users U ON C.user_id = U.id WHERE C.post_id = ?
	ORDER BY modified_at DESC LIMIT 10 OFFSET ` + strconv.Itoa((num-1)*10)
	rows, err = DB.Query(query, data.Id, postID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)

		data.Error(res, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Process comments
	for rows.Next() {
		var comment Comment
		var trash int
		var did string
		err = rows.Scan(&comment.Id, &comment.User.Id, &trash, &comment.Body, &comment.Created_at, &comment.Modified_at, &comment.UserName, &comment.Image, &comment.Like, &comment.Dislike, &did)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
		}
		comment.Duration = CalculDuration(comment.Created_at)
		if did != "" {
			comment.Did = true
			comment.Liked = (did == "1")
		}
		post.Comments = append(post.Comments, comment)
	}


	data.Posts = append(data.Posts, post)
	data.RenderPage("post.html", res)
}

// Function to delete a post
// func (data *Page) PostsDelete(res http.ResponseWriter, req *http.Request) {
// 	if req.Method != "GET" {
// 		data.Error(res, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	strID := req.URL.Path[len("/posts/delete/"):]
// 	postID, err := strconv.Atoi(strID)
// 	if err != nil {
// 		data.Error(res, http.StatusNotFound)
// 		return
// 	}

// 	// Check if the post belongs to the current user
// 	query := "SELECT image FROM posts WHERE id = ? AND user_id = ?"
// 	row := DB.QueryRow(query, postID, data.Id)
// 	var image string
// 	err = row.Scan(&image)
// 	if err != nil {
// 		data.Error(res, http.StatusInternalServerError)
// 		return
// 	}

// 	// Delete the post and its image
// 	query = "DELETE FROM posts WHERE id = ? AND user_id = ?"
// 	_, err = DB.Exec(query, postID, data.Id)
// 	if err != nil {
// 		data.Error(res, http.StatusInternalServerError)
// 		return
// 	}
// 	if image != "" {
// 		os.Remove("images/posts/" + image)
// 	}
// 	http.Redirect(res, req, "/", http.StatusFound)
// }

// Function to update a post
// func (data *Page) PostsUpdate(res http.ResponseWriter, req *http.Request) {
// 	if req.Method != "POST" {
// 		data.Error(res, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	postID, title, body, oldImage := req.FormValue("post_id"), strings.TrimSpace(req.FormValue("title")), strings.TrimSpace(req.FormValue("body")), req.FormValue("oldimage")
// 	if postID == "" || title == "" || len(title) > 400 || body == "" || len(body) > 5000 {
// 		data.Error(res, http.StatusBadRequest)
// 		return
// 	}

// 	// Process image if exists
// 	file, fileheader, err := req.FormFile("image")
// 	var image string
// 	if err == nil {
// 		image = data.HandleImage("posts", file, fileheader)
// 	}

// 	// Check if the user is the owner of the post
// 	query := "SELECT user_id FROM posts WHERE id = ?"
// 	row := DB.QueryRow(query, postID)
// 	var userID int
// 	err = row.Scan(&userID)
// 	if err != nil || userID != data.Id {
// 		data.Error(res, http.StatusBadRequest)
// 		return
// 	}

// 	// Handle image updates and remove the old image if necessary
// 	if image != "" || oldImage == "oldimage" {
// 		query = "SELECT image FROM posts WHERE id = ? AND user_id = ?"
// 		row = DB.QueryRow(query, postID, data.Id)
// 		var existingImage string
// 		err := row.Scan(&existingImage)
// 		if err != nil {
// 			data.Error(res, http.StatusInternalServerError)
// 			return
// 		}

// 		// Update the post in the database
// 		updateQuery := "UPDATE posts SET title = ?, body = ?, modified_at = ?, image = ? WHERE id = ?"
// 		_, err = DB.Exec(updateQuery, title, body, time.Now().Unix(), image, postID)
// 		if err != nil {
// 			data.Error(res, http.StatusInternalServerError)
// 			return
// 		}
// 		if existingImage != "" {
// 			os.Remove("images/posts/" + existingImage)
// 		}
// 	} else {
// 		// Update the post in the database
// 		updateQuery := "UPDATE posts SET title = ?, body = ?, modified_at = ? WHERE id = ?"
// 		_, err = DB.Exec(updateQuery, title, body, time.Now().Unix(), postID)
// 		if err != nil {
// 			data.Error(res, http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	http.Redirect(res, req, "/posts/"+postID, http.StatusFound)
// }
