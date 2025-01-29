package forum

import (
	"net/http"
	"strings"
)

// FilterCategories filters posts by selected categories.
func (data *Page) FilterCategories(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := req.ParseForm(); err != nil {
		data.Error(res, http.StatusBadRequest)
		return
	}

	categories := req.Form["categories"]
	if len(categories) == 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Use placeholders to prevent SQL injection
	query := `
		FROM posts_categories PC
		JOIN posts P ON PC.post_id = P.id
		JOIN users U ON P.user_id = U.id
		WHERE PC.category_id IN (` + placeholders(len(categories)) + `)
	`
	data.LoadPosts(req, res, query, stringSliceToInterface(categories)...)
}

// LikedPosts loads posts liked by the logged-in user.
func (data *Page) LikedPosts(res http.ResponseWriter, req *http.Request) {
	query := `
		FROM posts P
		JOIN users U ON P.user_id = U.id
		JOIN likes L ON P.id = L.post_id
		WHERE L.like = 1 AND L.user_id = ?
	`
	data.LoadPosts(req, res, query, data.Id)
}

// CreatedPosts loads posts created by the logged-in user.
func (data *Page) CreatedPosts(res http.ResponseWriter, req *http.Request) {
	query := `
		FROM posts P
		JOIN users U ON P.user_id = U.id
		WHERE U.id = ?
	`
	data.LoadPosts(req, res, query, data.Id)
}

// Helpers

// placeholders generates a comma-separated list of placeholders for SQL queries.
func placeholders(n int) string {
	return strings.Repeat("?,", n)[:2*n-1]
}

// stringSliceToInterface converts a slice of strings to a slice of interfaces.
func stringSliceToInterface(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
