package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CalculDuration calculates the duration since a post was created
func CalculDuration(Created_at int) string {
	dur := int(time.Now().Unix() - int64(Created_at))
	switch {
	case dur < 60:
		return fmt.Sprintf("%d seconds", dur)
	case dur < 3600:
		return fmt.Sprintf("%d minutes", dur/60)
	case dur < 86400:
		return fmt.Sprintf("%d hours", dur/3600)
	case dur < 2592000:
		return fmt.Sprintf("%d days", dur/86400)
	case dur < 31104000:
		return fmt.Sprintf("%d months", dur/2592000)
	default:
		return fmt.Sprintf("%d years", dur/31104000)
	}
}

// makeUrl generates pagination URLs
func makeUrl(pagenum int, basePath, oldQuery string) string {
	newQuery := []string{}
	if oldQuery != "" {
		newQuery = strings.Split(oldQuery, "&")
	}

	found := false
	for i, param := range newQuery {
		if strings.HasPrefix(param, "page=") {
			found = true
			newQuery[i] = fmt.Sprintf("page=%d", pagenum)
			break
		}
	}
	if !found {
		newQuery = append(newQuery, fmt.Sprintf("page=%d", pagenum))
	}
	return basePath + "?" + strings.Join(newQuery, "&")
}

// LoadPosts loads and processes posts based on a query and optional arguments
func (data *Page) LoadPosts(req *http.Request, res http.ResponseWriter, halfquery string, args ...any) {
	page := req.FormValue("page")
	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		data.Error(res, http.StatusBadRequest)
		return
	}

	basePath := req.URL.Path
	oldQuery := req.URL.RawQuery

	// Set pagination URLs
	data.Previous = makeUrl(pageNum-1, basePath, oldQuery)
	if pageNum == 1 {
		data.Previous = "0"
	}
	data.Next = makeUrl(pageNum+1, basePath, oldQuery)

	// Count total posts for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (SELECT DISTINCT P.* %s)", halfquery)
	var totalPosts int
	if err := DB.QueryRow(countQuery, args...).Scan(&totalPosts); err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	if totalPosts%10 != 0 {
		totalPosts += 10
	}
	data.Last = makeUrl(totalPosts/10, basePath, oldQuery)
	data.Current = pageNum != totalPosts/10 && totalPosts != 0

	// Load posts
	query := fmt.Sprintf(`
		SELECT DISTINCT P.*, U.username, U.image,
		(SELECT COUNT(*) FROM likes WHERE post_id = P.id AND like = 1) AS post_likes,
		(SELECT COUNT(*) FROM likes WHERE post_id = P.id AND like = 0) AS post_dislikes,
		COALESCE((SELECT like FROM likes WHERE post_id = P.id AND user_id = ?), "") AS did
		%s ORDER BY modified_at DESC LIMIT 10 OFFSET ?
	`, halfquery)
	args = append([]any{data.Id}, args...)
	rows, err := DB.Query(query, append(args, (pageNum-1)*10)...)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	// Process posts
	for rows.Next() {
		post := Post{}
		str := ""
		err := rows.Scan(&post.Id, &post.User.Id, &post.Title, &post.Body, &post.Image, &post.Created_at, &post.Modified_at, &post.UserName, &post.User.Image, &post.Like, &post.Dislike, &str)
		if err != nil {
			fmt.Println(err)
		}
		post.Duration = CalculDuration(post.Created_at)
		if str != "" {
			post.Did = true
			post.Liked = str == "1"
		}
		data.Posts = append(data.Posts, post)
	}
	rows.Close()

	for i, post := range data.Posts {
		query = "SELECT C.* FROM posts_categories PC JOIN categories C ON PC.category_id = C.id Where PC.post_id = ? "
		catrows, err := DB.Query(query, post.Id)
		if err != nil {
			data.Error(res, http.StatusBadRequest)
			return
		}
		for catrows.Next() {
			category := Categorie{}
			err = catrows.Scan(&category.Id, &category.CatName)
			if err != nil {
				fmt.Println(err)
			}
			data.Posts[i].Categories = append(data.Posts[i].Categories, category)
		}
		catrows.Close()
		trash := 0
		query = `	SELECT
						C.*,
						U.username,
						U.image
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
		row := DB.QueryRow(query, post.Id)
		err = row.Scan(&data.Posts[i].TopComment.Id, &data.Posts[i].TopComment.User.Id, &trash, &data.Posts[i].TopComment.Body, &data.Posts[i].TopComment.Created_at, &data.Posts[i].TopComment.Modified_at, &data.Posts[i].TopComment.UserName, &data.Posts[i].TopComment.Image)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
		}
	}
	data.RenderPage("index.html", res)
}

// HandleHome handles requests to the home page
func (data *Page) HandleHome(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	query := "FROM posts P JOIN users U ON P.user_id = U.id"
	data.LoadPosts(req, res, query)
}
