package server

import (
	"forum/internal"
	"forum/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppContext) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		s.ErrorHandler(w, http.StatusNotFound, "Not Found")
		return
	}
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, http.StatusNotAcceptable, "That email already occupied. Try another.")
		return
	}

	ok := s.alreadyLogIn(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		var u models.User
		u.UserName = r.FormValue("uname")
		u.Email = r.FormValue("uemail")
		pass := r.FormValue("upass")
		passConfirm := r.FormValue("confirm-upass")

		if ok := internal.CheckName(u.UserName); !ok {
			s.ErrorHandler(w, http.StatusForbidden, "Please, be sure to write valid name")
			return
		}
		if ok := internal.CheckEmail(u.Email); !ok {
			s.ErrorHandler(w, http.StatusForbidden, "Please, be sure to write valid email")
			return
		}

		if pass != passConfirm || len(pass) > 100 {
			s.ErrorHandler(w, http.StatusNotAcceptable, "Password doesn't match / more than 100 char")
			return
		}

		if len(pass) < 5 {
			s.ErrorHandler(w, http.StatusNotAcceptable, "Password should be more than 5 symbols")
			return
		}

		var err error
		u.Password, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
		CheckErr(err)

		_, err = s.Sqlite3.InsertUser(&u)
		if err != nil {
			s.ErrorHandler(w, http.StatusNotAcceptable, "That email already occupied. Try another.")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

}

func (s *AppContext) alreadyLogIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	

	// getSession from db
	_, err = s.Sqlite3.GetUserID(cookie.Value)

	return err == nil
}
