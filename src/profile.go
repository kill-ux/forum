package forum

import (
	"net/http"
	"os"
)

func (data Page) UpdatePFP(res http.ResponseWriter, req *http.Request) {
	// DELETE OLD pfp
	query := "SELECT image FROM users WHERE id = ? "
	row := data.DB.QueryRow(query, data.Id)
	oldPFP := ""
	err := row.Scan(&oldPFP)
	if err != nil {
		data.Error(res, http.StatusInternalServerError)
		return
	}

	// fetching form data
	file, fileheader, err := req.FormFile("image")
	var image string
	if err == nil {
		image = data.HandelImage("pics", file, fileheader, res, req)
	}
	// UPDATE DB
	if image != "" {
		query = "UPDATE users SET image = ? WHERE id = ?"
		_, err = data.DB.Exec(query, image, data.Id)
		if err != nil {
			data.Error(res, http.StatusInternalServerError)
			return
		}
		if oldPFP != "profile.png" {
			os.Remove("images/pics/" + oldPFP)
		}
	}
	// REDIRECT
	http.Redirect(res, req, "/", http.StatusFound)
}
