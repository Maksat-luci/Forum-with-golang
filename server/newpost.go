package server

import (
	"forum/internal"
	"forum/pkg/models"
	"net/http"
)

func (s *AppContext) showNewPost(w http.ResponseWriter, r *http.Request) {

	if !s.alreadyLogIn(r) {
		s.ErrorHandler(w, http.StatusForbidden, "please log-in first")
		return
	}

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300 // 300 is session length

	// update session table last activity
	userID, err := s.Sqlite3.GetUserID(cookie.Value)
	if err != nil {
		s.ErrorLog.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	if r.URL.Path != "/newpost" {
		s.ErrorLog.Println(r.URL.Path)
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}

	switch r.Method {
	case http.MethodGet:
		categories, err := s.Sqlite3.GetAllCategories()
		if err != nil {
			s.ErrorLog.Println(err)
			s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = s.Template.ExecuteTemplate(w, "newpost.html", categories)
		if err != nil {
			s.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		s.newPost(w, r)
	default:
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")

	}

}

func (s *AppContext) newPost(w http.ResponseWriter, r *http.Request) {
	if ok := s.alreadyLogIn(r); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	cookie, _ := r.Cookie("session")
	userID, _ := s.Sqlite3.GetUserID(cookie.Value)

	if len(internal.RemoveSpace(r.FormValue("category"))) == 0 {
		s.ErrorLog.Println("emty post")
		s.ErrorHandler(w, http.StatusBadRequest, "Choose Category")
		return
	}
	var p models.Post
	p.UserID = userID
	p.Title = r.FormValue("title")
	p.Content = r.FormValue("post")

	ok := internal.CheckPost(&p)
	switch ok {
	case 1:
		s.ErrorHandler(w, http.StatusBadRequest, "Empty title / post")
		return
	case 2:
		s.ErrorHandler(w, http.StatusBadRequest, "Post's content should not be more 150 char")
		return
	case 3:
		s.ErrorHandler(w, http.StatusBadRequest, "Post's title should not be more 60 char")
		return
	}

	// putting recieved data into database
	postID, err := s.Sqlite3.InsertPost(&p)
	if err != nil {
		s.ErrorLog.Println(err)
		s.ErrorHandler(w, http.StatusBadRequest, "Post titles should be unique")
		return

	}
	postCategories, err := internal.IsValidCtgry(r)
	if err != nil {
		s.ErrorLog.Println(err)
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	err = s.Sqlite3.AddPostCategory(postCategories, postID)

	if err != nil {
		s.ErrorLog.Println(err)
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
