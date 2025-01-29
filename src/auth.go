package forum

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// Utility function to set cookies
func setErrorCookie(res http.ResponseWriter, message, path string, maxAge int) {
	http.SetCookie(res, &http.Cookie{
		Name:     "errors",
		Value:    message,
		Path:     path,
		MaxAge:   maxAge,
		// HttpOnly: true, // Secure the cookie, not accessible by JS
	})
}

// Render the login page
func (data Page) Login(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	data.RenderPage("login.html", res)
}

// Render the signup page
func (data Page) Signup(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}
	data.RenderPage("signup.html", res)
}

// Handle user login
func (data *Page) AuthLogin(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	email, pass := strings.ToLower(req.FormValue("email")), req.FormValue("password")

	// Validate input
	if len(email) > 400 || len(pass) < 8 || len(pass) > 400 {
		setErrorCookie(res, "Invalid credentials!", "/login", 2)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	// Fetch user from the database
	query := `SELECT * FROM users WHERE email = ? OR username = ?`
	var user User
	row := DB.QueryRow(query, email, email)
	err := row.Scan(&user.Id, &user.UserName, &user.Email, &user.Password, &user.Image, &user.Role, &user.Token, &user.Token_Exp, &user.Created_at)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) != nil {
		setErrorCookie(res, "Incorrect credentials!", "/login", 2)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	// Generate and set a new token
	ntoken, err := uuid.NewV1()
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   ntoken.String(),
		Expires: time.Now().Add(TOKEN_AGE),
		Path:    "/",
	})

	// Update token in the database
	query = `UPDATE users SET token = ?, token_exp = ? WHERE id = ?`
	if _, err := DB.Exec(query, ntoken, time.Now().Add(24*time.Hour).Unix(), user.Id); err != nil {
		setErrorCookie(res, "Unexpected error, try again!", "/login", 2)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	// Log the user in
	// data.User = user
	data.Log = 1
	setErrorCookie(res, "", "/login", -1) // Remove the error cookie
	http.Redirect(res, req, "/", http.StatusFound)
}

// Handle user signup
func (data Page) AuthSignup(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	name, email, pass, confirm := strings.ToLower(req.FormValue("name")), strings.ToLower(req.FormValue("email")), req.FormValue("password"), req.FormValue("confirm")

	// Validate input
	emailRg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRg := regexp.MustCompile(`^[a-zA-Z0-9-]{3,100}$`)
	if pass != confirm || !emailRg.MatchString(email) || len(email) > 200 || !nameRg.MatchString(name) ||
		len(name) > 400 || len(pass) < 8 || len(pass) > 400 {
		setErrorCookie(res, "Invalid credentials!", "/signup", 2)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}

	// Hash password
	password, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		setErrorCookie(res, "Unexpected error, try again!", "/signup", 2)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}

	// Insert user into the database
	query := `INSERT INTO users VALUES (NULL, ?, ?, ?, 'profile.png', 'user', NULL, 0, ?)`
	if _, err := DB.Exec(query, name, email, string(password), time.Now()); err != nil {
		setErrorCookie(res, "Invalid data, try again!", "/signup", 2)
		http.Redirect(res, req, "/signup", http.StatusFound)
		return
	}

	// Redirect to login page on success
	setErrorCookie(res, "", "/signup", -1) // Remove the error cookie
	http.Redirect(res, req, "/login", http.StatusFound)
}

// Verify user token
func (data *Page) VerifyToken(token string) bool {
	if _, err := uuid.FromString(token); err != nil {
		return false
	}

	query := `SELECT * FROM users WHERE token = ?`
	row := DB.QueryRow(query, token)
	err := row.Scan(&data.Id, &data.UserName, &data.Email, &data.Password, &data.Image, &data.Role, &data.Token, &data.Token_Exp, &data.Created_at)

	validDate := int(time.Now().Unix()) < data.Token_Exp
	return err == nil && validDate
}

// Handle user logout
func (data *Page) AuthLogout(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		data.Error(res, http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})

	data.User = User{}
	http.Redirect(res, req, "/login", http.StatusFound)
}
