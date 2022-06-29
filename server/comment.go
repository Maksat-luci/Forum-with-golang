package server

import (
	"net/http"
	"strconv"
	"strings"
)

func (s *AppContext) comment(w http.ResponseWriter, r *http.Request) {
	if !s.alreadyLogIn(r) {
		s.ErrorHandler(w, http.StatusForbidden, "please, log-in first")
		return
	}

	if r.Method != http.MethodPost {
		s.ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	cookie, _ := r.Cookie("session")
	cookie.MaxAge = 300

	// update session table last activity
	userID, err := s.Sqlite3.GetUserID(cookie.Value)
	if err != nil {
		s.ErrorHandler(w, 500, "Internal Server Error")
		return
	}
	if s.Sqlite3.HasSession(userID) {
		s.Sqlite3.UpdateSession(userID)
	}

	postID, err := strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		s.ErrorHandler(w, http.StatusBadRequest, "Bad Request")
		return
	}
	content := r.FormValue("content")
	if len(strings.Trim(content, " ")) == 0 {
		s.ErrorLog.Println("emty comment")
		s.ErrorHandler(w, http.StatusForbidden, "Empty comment forbidden")
		return
	} else if len(content) > 140 {
		s.ErrorHandler(w, http.StatusForbidden, "Empty should not be more 140 char")
		return
	}
	err = s.Sqlite3.AddComment(userID, postID, content)
	if err != nil {
		s.ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	url := "/post/" + strconv.Itoa(postID)
	http.Redirect(w, r, url, http.StatusSeeOther)

}
