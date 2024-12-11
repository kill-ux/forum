package forum

import (
	"html/template"
	"net/http"
)

// render page of html
func (data Page) RenderPage(file string, res http.ResponseWriter) {
	// Parse all templates
	temp, err := template.ParseFiles(
		"views/header.html",
		"views/footer.html",
		"views/"+file,
	)
	if err != nil {
		if file == "error.html" {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		} else {
			data.Error(res, http.StatusInternalServerError)
			return
		}
	}
	err = temp.ExecuteTemplate(res, file, data)
	if err != nil {
		if file == "error.html" {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		} else {
			data.Error(res, http.StatusInternalServerError)
			return
		}
	}
}

func (data Page) Error(res http.ResponseWriter, code int) {
	res.WriteHeader(code)
	data.PageError.Code = code
	data.PageError.Msg = http.StatusText(code)
	data.RenderPage("error.html", res)
}
