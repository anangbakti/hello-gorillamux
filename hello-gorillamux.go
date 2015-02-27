package hello

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Book struct {
	Title  string
	Author string
}

func init() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/posts", PostsOnlyIndexHandler)
	r.HandleFunc("/posts/{id}", PostsIndexHandler)
	r.HandleFunc("/view1/", View1Handler)
	http.Handle("/", r)
}

var tpl = template.Must(template.ParseGlob("tmpl/*.html"))

//contoh handler call tempate html
func View1Handler(rw http.ResponseWriter, r *http.Request) {
	book := Book{"Building Web Apps with Go", "Jeremy Saenz"}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(rw, "indexPage", book); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Homeku")
}

func PostsOnlyIndexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, fmt.Sprintf("at %v posts", "hello"))
}

func PostsIndexHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintln(rw, fmt.Sprintf("at %v, %s", "hello", id))

}
