package forum

import (
	"net/http"
	"os"
	"strings"
)

func (data *Page) ResetUser(req *http.Request) {
	emailerror, err := req.Cookie("errors")
	if err == nil {
		data.User = User{ErrorEmail: emailerror.Value}
	} else {
		data.User = User{}
	}
}

func checkSuff(path string) bool {
	// arr := []string{"store", "update", "likes"}
	arr := []string{"store", "likes"}
	for _, v := range arr {
		if strings.HasSuffix(path, v) {
			return false
		}
	}
	return true
}

func Routers(res http.ResponseWriter, req *http.Request) {
	data := &Page{}
	data.FillCategories()
	data.ResetUser(req)
	token, err := req.Cookie("token")

	path := req.URL.Path[1:]
	if path == "" || (strings.HasPrefix(path, "posts/") && checkSuff(path)) || path == "filter/categories" || strings.HasPrefix(path, "css/") || strings.HasPrefix(path, "js/") || strings.HasPrefix(path, "images/") {
		if err == nil && data.VerifyToken(token.Value) {
			data.Log = 1
		}
	} else if path != "login" && path != "auth/login" && path != "signup" && path != "auth/signup" {
		if err != nil || !data.VerifyToken(token.Value) {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}
		data.Log = 1
	} else {
		if err == nil && data.VerifyToken(token.Value) {
			data.Log = 1
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}
	}

	switch path {
	case "":
		data.HandleHome(res, req)
	case "login":
		data.Login(res, req)
	case "auth/login":
		data.AuthLogin(res, req)
	case "signup":
		data.Signup(res, req)
	case "auth/signup":
		data.AuthSignup(res, req)
	case "auth/logout":
		data.AuthLogout(res, req)
	case "posts/store":
		data.PostsStore(res, req)
	case "comments/store":
		data.CommentsStore(res, req)
	case "comments/likes":
		data.CommentsLike(res, req)
	case "posts/likes":
		data.PostsLike(res, req)
	// case "posts/update":
	// 	data.PostsUpdate(res, req)
	case "filter/categories":
		data.FilterCategories(res, req)
	case "filter/created":
		data.CreatedPosts(res, req)
	case "filter/liked":
		data.LikedPosts(res, req)
	// case "pfp/update":
	// 	data.UpdatePFP(res, req)
	default:
		// if strings.HasPrefix(path, "posts/delete/") {
		// 	data.PostsDelete(res, req)
		// } else
		if strings.HasPrefix(path, "css/") {
			_, err := os.ReadFile(req.URL.Path[1:])
			if err != nil {
				data.Error(res, http.StatusNotFound)
				return
			}
			http.StripPrefix("/css/", http.FileServer(http.Dir("css"))).ServeHTTP(res, req)
		} else if strings.HasPrefix(path, "js/") {
			_, err := os.ReadFile(req.URL.Path[1:])
			if err != nil {
				data.Error(res, http.StatusInternalServerError)
				return
			}
			http.StripPrefix("/js/", http.FileServer(http.Dir("js"))).ServeHTTP(res, req)
		} else if strings.HasPrefix(path, "images/") {
			_, err := os.ReadFile(req.URL.Path[1:])
			if err != nil {
				data.Error(res, http.StatusInternalServerError)
				return
			}
			http.StripPrefix("/images/", http.FileServer(http.Dir("images"))).ServeHTTP(res, req)
		} else if strings.HasPrefix(path, "posts/") && strings.HasSuffix(path, "/comments") {
			data.CommentsInfo(res, req)
		} else if strings.HasPrefix(path, "posts/") {
			data.PostsInfo(res, req)
		} else {
			data.Error(res, http.StatusNotFound)
			return
		}
	}
}
