package ushow

// import (
// 	"html/template"
// 	"net/http"
// )

// //Compile templates on start
// var templates = template.Must(template.ParseFiles(
// 	"../web/upload.html",
// 	"../web/shared_session.html"))

// //Display the named template
// func display(w http.ResponseWriter, tmpl string, data interface{}) {
// 	templates.ExecuteTemplate(w, tmpl+".html", data)
// }

// func DisplaySharedSession(w http.ResponseWriter, r *http.Request) {
// 	display(w, "shared_session", "test")
// }
