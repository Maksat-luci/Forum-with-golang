package server

import (
	"forum/pkg/models/sqlite3"
	"html/template"
	"log"
)

type AppContext struct {
	Sqlite3  *sqlite3.Database
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Template *template.Template
	
}

func NewAppContext(db *sqlite3.Database, InfoLog, ErrorLog *log.Logger, tmpl *template.Template) *AppContext {
	return &AppContext{
		Sqlite3:  db,
		InfoLog:  InfoLog,
		ErrorLog: ErrorLog,
		Template: tmpl,
	}
}
