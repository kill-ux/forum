package forum

import (
	"net/http"
)

func (data *Page) ResetUser(req *http.Request) {
	emailerror, err := req.Cookie("errors")
	if err == nil {
		data.User = User{ErrorEmail: emailerror.Value}
	} else {
		data.User = User{}
	}
}

func (data *Page) Routers(res http.ResponseWriter, req *http.Request) {
	data.ResetUser(req)
	data.Posts = data.Posts[:0]

	token, err := req.Cookie("token")
	path := req.URL.Path[1:]
	if path == "" || (len(path) > len("posts/") && path[:len("posts/")] == "posts/" && path[len("posts/"):] != "store" && path[len("posts/"):] != "likes") || path == "filter/categories" {
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
	case "posts/update":
		data.PostsUpdate(res, req)
	case "filter/categories":
		data.FilterCategories(res, req)
	case "filter/created":
		data.CreatedPosts(res, req)
	case "filter/liked":
		data.LikedPosts(res, req)
	case "pfp/update":
		data.UpdatePFP(res, req)
	default:
		if len(path) > len("posts/delete/") && path[:len("posts/delete/")] == "posts/delete/" {
			data.PostsDelete(res, req)
		} else if len(path) > len("posts/") && path[:len("posts/")] == "posts/" {
			data.PostsInfo(res, req)
		} else {
			data.Error(res, http.StatusNotFound)
			return
		}
	}
}
