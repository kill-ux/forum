package forum

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleLike performs the common logic for managing likes
func handleLike(res http.ResponseWriter, req *http.Request, likeType, idKey string, data Page) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := req.FormValue(idKey)
	like := req.FormValue("like")
	if like != "0" && like != "1" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	var additionalCondition string
	if likeType == "comment" {
		additionalCondition = "AND post_id ISNULL"
	} else if likeType == "post" {
		additionalCondition = "AND comment_id ISNULL"
	}

	// Check if a like/dislike already exists
	query := "SELECT like FROM likes WHERE user_id = ? AND " + idKey + " = ? " + additionalCondition
	row := DB.QueryRow(query, data.Id, id)

	var liked bool
	err := row.Scan(&liked)
	if err != nil { // No existing like/dislike, insert a new one
		query = "INSERT INTO likes (like, user_id, " + idKey + ") VALUES (?, ?, ?)"
		DB.Exec(query, like, data.Id, id)
	} else {
		// If like status matches, remove the like; otherwise, update it
		if (like == "1" && liked) || (like == "0" && !liked) {
			query = "DELETE FROM likes WHERE user_id = ? AND " + idKey + " = ? " + additionalCondition
			DB.Exec(query, data.Id, id)
		} else {
			query = "UPDATE likes SET like = ? WHERE user_id = ? AND " + idKey + " = ? " + additionalCondition
			DB.Exec(query, like, data.Id, id)
		}
	}

	var Response struct {
		Likes    int
		Dislikes int
		Did      bool
		Like     bool
	}
	did := ""
	query = fmt.Sprintf(`SELECT (SELECT COUNT(*) FROM likes WHERE %s = ? AND like = 1) as comment_likes,
		(SELECT COUNT(*) FROM likes WHERE %s = ? AND like = 0) as comment_dislikes,
		COALESCE((SELECT like FROM likes WHERE %s = ? AND user_id = ?), "") as did`, idKey, idKey, idKey)
	if err = DB.QueryRow(query, id, id, id, data.Id).Scan(&Response.Likes, &Response.Dislikes, &did); err != nil {
		fmt.Println(err)
	}
	if did != "" {
		Response.Did = true
		Response.Like = (did == "1")
	}

	json.NewEncoder(res).Encode(Response)

	// Use the getPathFromReferer function to properly handle the redirect URL
	// refURL := getPathFromReferer(req)
	// http.Redirect(res, req, refURL+"#"+likeType+id, http.StatusFound)
}

// CommentsLike handles likes and dislikes for comments
func (data Page) CommentsLike(res http.ResponseWriter, req *http.Request) {
	handleLike(res, req, "comment", "comment_id", data)
}

// PostsLike handles likes and dislikes for posts
func (data Page) PostsLike(res http.ResponseWriter, req *http.Request) {
	handleLike(res, req, "post", "post_id", data)
}
