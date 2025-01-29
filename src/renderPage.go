package forum

import (
	"html/template"
	"net/http"
)

// RenderPage renders the HTML page with the provided file and data
func (data Page) RenderPage(file string, res http.ResponseWriter) {
	// Handle error email if set
	if data.User.ErrorEmail != "" {
		res.WriteHeader(http.StatusBadRequest)
	}

	// Parse the templates
	temp, err := template.ParseFiles(
		"views/header.html",
		"views/footer.html",
		"views/"+file,
	)
	if err != nil {
		data.handleTemplateError(err, res, file)
		return
	}

	// Execute the template with data
	err = temp.ExecuteTemplate(res, file, data)
	if err != nil {
		data.handleTemplateError(err, res, file)
	}
}

// handleTemplateError handles errors related to template rendering
func (data Page) handleTemplateError(err error, res http.ResponseWriter, file string) {
	if file == "error.html" {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	} else {
		data.Error(res, http.StatusInternalServerError)
	}
}

// Error renders an error page with the specified status code
func (data Page) Error(res http.ResponseWriter, code int) {
	res.WriteHeader(code)
	data.PageError.Code = code
	data.PageError.Msg = http.StatusText(code)
	data.RenderPage("error.html", res)
}
