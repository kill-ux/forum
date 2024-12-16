package forum

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func (data Page) Login(res http.ResponseWriter, req *http.Request) {
	data.RenderPage("login.html", res)
}

func (data Page) Signup(res http.ResponseWriter, req *http.Request) {
	data.RenderPage("signup.html", res)
}

func (data *Page) AuthLogin(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	errcookie := http.Cookie{
		Name:     "errors",
		Path:     "/login",
		MaxAge:   2,
		HttpOnly: true, // Secure the cookie, not accessible by JS
	}
	email, pass := req.FormValue("email"), req.FormValue("password")
	emailRg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRg.MatchString(email) || len(email) > 400 || len(pass) < 8 || len(pass) > 400 {
		errcookie.Value = "invalide credentials!"
		http.SetCookie(res, &errcookie)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	query := `SELECT * FROM users WHERE email = ? OR username = ? `
	var user User
	row := data.DB.QueryRow(query, email, email)
	err := row.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Image, &user.Role, &user.Token, &user.Token_Exp, &user.Created_at)
	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil || err1 != nil {
		errcookie.Value = "incorrect credentials!"
		http.SetCookie(res, &errcookie)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	data.User = user
	ntoken, err := uuid.NewV1()
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	cookies := &http.Cookie{
		Name:    "token",
		Value:   ntoken.String(),
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(res, cookies)
	query = `UPDATE users SET token = ? , token_exp = ? WHERE id = ?`
	_, err = data.DB.Exec(query, ntoken, time.Now().Add(24*time.Hour).Unix(), data.Id)
	if err != nil {
		errcookie.Value = "unexpected error try again!"
		http.SetCookie(res, &errcookie)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	data.Log = 1
	// remove errors if logged in succefully
	errcookie.MaxAge = -1
	http.SetCookie(res, &errcookie)
	http.Redirect(res, req, "/", http.StatusFound)
}

func (data Page) AuthSignup(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	cookie := http.Cookie{
		Name:     "errors",
		Path:     "/signup",
		MaxAge:   2,
		HttpOnly: true, // Secure the cookie, not accessible by JS
	}
	name, email, pass, confirm := req.FormValue("name"), req.FormValue("email"), req.FormValue("password"), req.FormValue("confirm")
	emailRg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRg := regexp.MustCompile(`^[a-zA-Z0-9-]{3,100}$`)
	if pass != confirm || !emailRg.MatchString(email) || len(email) > 200 || !nameRg.MatchString(name) ||
		len(name) > 400 || len(pass) < 8 || len(pass) > 400 || len(confirm) < 8 || len(confirm) > 400 {
		cookie.Value = "invalide credentials!"
		http.SetCookie(res, &cookie)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		cookie.Value = "unexpected error try again!"
		http.SetCookie(res, &cookie)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}
	pass = string(password)
	query := `INSERT INTO users Values (NULL,?,?,?,'profile.png','user',NULL,0,?)`
	_, err = data.DB.Exec(query, name, email, pass, time.Now())
	if err != nil {
		cookie.Value = "unvalid data, try again!"
		http.SetCookie(res, &cookie)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(res, &cookie)
	http.Redirect(res, req, "/login", http.StatusFound)
}

func (data *Page) VerifyToken(token string) bool {
	_, err := uuid.FromString(token)
	query := `SELECT * FROM users WHERE token = ? `
	row := data.DB.QueryRow(query, token)
	err1 := row.Scan(&data.Id, &data.UserName, &data.Email, &data.Password, &data.Image, &data.Role, &data.Token, &data.Token_Exp, &data.Created_at)
	dur := int(time.Now().Add(TOKEN_AGE).Unix()) - data.Token_Exp
	return err == nil && err1 == nil && dur < int(TOKEN_AGE)
}

func (data *Page) AuthLogout(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	cookies := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	}
	http.SetCookie(res, cookies)
	data.User = User{}
	http.Redirect(res, req, "/login", http.StatusFound)
}
