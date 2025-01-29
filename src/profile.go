package forum

// import (
// 	"net/http"
// 	"os"
// )

// func (data Page) UpdatePFP(res http.ResponseWriter, req *http.Request) {
// 	// Retrieve the current profile picture (PFP)
// 	query := "SELECT image FROM users WHERE id = ?"
// 	row := DB.QueryRow(query, data.Id)
// 	var oldPFP string
// 	err := row.Scan(&oldPFP)
// 	if err != nil {
// 		data.Error(res, http.StatusInternalServerError)
// 		return
// 	}

// 	// Handle the image upload
// 	file, fileheader, err := req.FormFile("image")
// 	var image string
// 	if err == nil {
// 		image = data.HandleImage("pics", file, fileheader)
// 	} else {
// 		data.Error(res, http.StatusInternalServerError)
// 		return
// 	}

// 	// If a new image is provided, update the database and delete the old one if necessary
// 	if image != "" {
// 		// Update the database with the new image
// 		query = "UPDATE users SET image = ? WHERE id = ?"
// 		_, err = DB.Exec(query, image, data.Id)
// 		if err != nil {
// 			data.Error(res, http.StatusInternalServerError)
// 			return
// 		}

// 		// Remove the old profile picture if it’s not the default one
// 		if oldPFP != "profile.png" && oldPFP != "" {
// 			err := os.Remove("images/pics/" + oldPFP)
// 			if err != nil {
// 				data.Error(res, http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	}

// 	// Redirect the user back to the home page or another suitable location
// 	http.Redirect(res, req, "/", http.StatusFound)
// }
