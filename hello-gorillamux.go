package hello

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
)

type Book struct {
	Title  string
	Author string
}

type Location struct {
	Loc_Name string
}

func init() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/posts", PostsOnlyIndexHandler)
	r.HandleFunc("/posts/{id}", PostsIndexHandler)
	r.HandleFunc("/view1/", View1Handler)
	r.HandleFunc("/entryapi/", EntryApiHandler)
	r.HandleFunc("/saveloc", SaveLocHandler)
	r.HandleFunc("/viewentry/", ViewEntryApiHandler)
	r.HandleFunc("/locations/", AllLocHandler)
	http.Handle("/", r)
}

var tpl = template.Must(template.ParseGlob("tmpl/*.html"))

func AllLocHandler(rw http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Location").Ancestor(LocationKey(c)).Limit(10)
	var locs []*Location
	if _, err := q.GetAll(c, &locs); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(locs); err != nil {
		panic(err)
	}
}

//contoh handler call tempate html
func View1Handler(rw http.ResponseWriter, r *http.Request) {
	book := Book{"Building Web Apps with Go", "Jeremy Saenz"}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(rw, "indexPage", book); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func EntryApiHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(rw, "entryApiPage", nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func ViewEntryApiHandler(rw http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Location").Ancestor(LocationKey(c)).Limit(10)
	var locs []*Location
	if _, err := q.GetAll(c, &locs); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(rw, "viewEntryApiPage", locs); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func SaveLocHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(rw, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	c := appengine.NewContext(r)
	l := &Location{
		Loc_Name: r.FormValue("loc_name"),
	}
	//fmt.Fprintln(rw, r.FormValue("loc_name"))
	key := datastore.NewIncompleteKey(c, "Location", LocationKey(c))
	if _, err := datastore.Put(c, key, l); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/entryapi/", http.StatusSeeOther)
}

// guestbookKey returns the key used for all guestbook entries.
func LocationKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Location", "default_loc", 0, nil)
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
