package server

import (
	"forum/pkg/models"
	"net/http"
)

func (s *AppContext)index(w http.ResponseWriter, r *http.Request) {
	var login int
	if s.alreadyLogIn(r) {
		login = 1
	}
	s.InfoLog.Printf("Call page: %v\t\tHandleFunc: index", r.URL.Path)
	if r.URL.Path != "/" {
		s.ErrorLog.Println("r.URL.Path: ", r.URL.Path)
		s.ErrorHandler(w, http.StatusNotFound, "Not Found")
		return
	}

	if r.Method != http.MethodGet {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// new function
	allPosts, err := s.Sqlite3.GetAllPosts()
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// categories
	categories, err := s.Sqlite3.GetAllCategories()
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := struct {
		AllPosts   *[]models.Post
		Categories map[string]int
		Auth       int
	}{allPosts, categories, login}

	err = s.Template.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		s.InfoLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
