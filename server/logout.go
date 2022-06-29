package server

import (
	"net/http"
)

func (s *AppContext) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		s.ErrorLog.Println(r.URL.Path)
		s.ErrorHandler(w, http.StatusNotFound, "Not Found")
		return
	}

	ok := s.alreadyLogIn(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		s.ErrorLog.Println("cookie doesn't exist")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := s.Sqlite3.GetUserID(cookie.Value)
	if err != nil {
		s.ErrorLog.Println("cookie doesn't exist")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	s.Sqlite3.DeleteSession(userID)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
