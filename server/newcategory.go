package server

import "net/http"

func (s *AppContext) newCategory(w http.ResponseWriter, r *http.Request) {
	s.InfoLog.Printf("Call page: %v\t\tHandleFunc: newCategory", r.URL.Path)
	if !s.alreadyLogIn(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}
	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	userID, _ := s.Sqlite3.GetUserID(cookie.Value)

	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("category")

		err := s.Sqlite3.InsertCategory(title)
		if err != nil {
			s.ErrorLog.Println(err)
			s.ErrorHandler(w, 500, "Internal Server Error")
			return
		}
		 
		s.InfoLog.Println("Added new category")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := s.Template.ExecuteTemplate(w, "newcategory.html", nil)
	if err != nil {
		s.ErrorLog.Println(err)
		s.ErrorHandler(w, 500, "Internal Server Error")
		return
	}
}
