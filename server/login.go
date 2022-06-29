package server

import (
	"fmt"
	"forum/internal"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		s.ErrorHandler(w, http.StatusNotFound, "Not Found")
		return
	}

	ok := s.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		err := s.Template.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		s.loginPost(w, r)
	default:
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

}

func (s *AppContext) loginPost(w http.ResponseWriter, r *http.Request) {
	var clientEmail, clientPass string

	clientEmail = r.FormValue("uemail")
	clientPass = r.FormValue("upass")

	if len(clientEmail) > 320 || len(clientPass) > 100 || !internal.CheckEmail(clientEmail) {
		errorMsg := struct {
			Msg string
		}{
			"incorrect login",
		}
		// 401 unauthorised
		w.WriteHeader(http.StatusUnauthorized)
		s.Template.ExecuteTemplate(w, "login.html", errorMsg)
		return
	}
	// email check in db
	u, err := s.Sqlite3.GetUser(clientEmail)

	if err != nil {
		s.ErrorLog.Println(err)
		errorMsg := struct {
			Msg string
		}{
			"login doesn't exist",
		}
		// 403 unauthorised
		w.WriteHeader(http.StatusBadRequest)
		err = s.Template.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			s.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(clientPass))
	if err != nil {
		s.ErrorLog.Println(err)
		errorMsg := struct {
			Msg string
		}{
			"incorrect password",
		}
		w.WriteHeader(401)
		err := s.Template.ExecuteTemplate(w, "login.html", errorMsg)
		if err != nil {
			s.ErrorLog.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	if s.Sqlite3.HasSession(u.UserID) {
		s.Sqlite3.DeleteSession(u.UserID)
	}
	//create new function
	sID := internal.SetCookie(w)

	err = s.Sqlite3.InsertSession(u.UserID, sID.String())
	if err != nil {
		s.ErrorLog.Println(err)
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("========== Logged-in successfully ==========")

}
