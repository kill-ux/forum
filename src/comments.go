package forum

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// CommentsStore handles the creation of comments for a post.
func (data Page) CommentsStore(res http.ResponseWriter, req *http.Request) {
	// Ensure the request method is POST
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	// Extract form values
	postID := req.FormValue("post_id")
	body := req.FormValue("body")
	jsBollen := req.FormValue("js")

	// Check if the user is on cooldown
	Mux.Lock()
	_, isOnCooldown := Cach[data.Id]
	Mux.Unlock()

	var ResponseError struct {
		Message string
		Code    int
	}

	if data.ErrorEmail != "" {
		if jsBollen == "" {
			data.Error(res, http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			ResponseError.Code = http.StatusBadRequest
			ResponseError.Message = data.ErrorEmail
			json.NewEncoder(res).Encode(ResponseError)
		}
		return
	}

	if isOnCooldown {
		if jsBollen == "" {
			setErrorCookie(res, "You are on cooldown! Wait for a bit and try again ^_^.", getPathFromReferer(req), 5)
			http.Redirect(res, req, "/posts/"+postID, http.StatusFound)
		} else {
			setErrorCookie(res, "You are on cooldown! Wait for a bit and try again ^_^.", getPathFromReferer(req), 5)
			res.WriteHeader(http.StatusBadRequest)
			ResponseError.Code = http.StatusBadRequest
			ResponseError.Message = "You are on cooldown! Wait for a bit and try again ^_^."
			json.NewEncoder(res).Encode(ResponseError)
		}
		return
	}

	// Validate the comment body length
	if len(body) < 1 || len(body) > 5000 {
		if jsBollen == "" {
			data.Error(res, http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			ResponseError.Code = http.StatusBadRequest
			ResponseError.Message = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(res).Encode(ResponseError)
		}
		return
	}

	// Insert the comment into the database
	query := "INSERT INTO comments VALUES (NULL,?,?,?,?,?)"
	_, err := DB.Exec(query, data.Id, postID, body, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		if jsBollen == "" {
			data.Error(res, http.StatusInternalServerError)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			ResponseError.Code = http.StatusInternalServerError
			ResponseError.Message = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(ResponseError)
		}
		return
	}

	// Add the user to the cooldown cache
	Mux.Lock()
	Cach[data.Id] = time.Now().Unix()
	Mux.Unlock()

	// Redirect the user back to the post
	if jsBollen == "" {
		http.Redirect(res, req, "/posts/"+postID, http.StatusFound)
	}
}

// getPathFromReferer extracts the path from the Referer header.
func getPathFromReferer(req *http.Request) string {
	referer := req.Referer()
	if referer == "" {
		return "/"
	}

	// Find the first '/' after the domain and return the path
	if idx := strings.Index(referer, "://"); idx != -1 {
		if pathIdx := strings.Index(referer[idx+3:], "/"); pathIdx != -1 {
			return referer[idx+3+pathIdx:]
		}
	}
	return "/"
}
