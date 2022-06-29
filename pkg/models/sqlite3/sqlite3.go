package sqlite3

import (
	"database/sql"
	"log"
)

type Database struct {
	SqlDb *sql.DB
}

func ConnectDb(driverName string, SqlDbName string) (*Database, error) {
	SqlDb, err := sql.Open(driverName, SqlDbName)
	if err != nil {
		return nil, err
	}
	if err = SqlDb.Ping(); err != nil {
		return nil, err
	}
	return &Database{SqlDb}, nil
}

func (c *Database) CreatePeopleTable()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "people" (
		"user_id"	INTEGER NOT NULL,
		"email"	VARCHAR(320) NOT NULL UNIQUE,
		"username"	VARCHAR(255) NOT NULL,
		"password"	BLOB NOT NULL,
		"time_creation" DATETIME,
		PRIMARY KEY("user_id" AUTOINCREMENT)
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	

}
func (c *Database) CreateSessionTable()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "sessions" (
		"user_id"	INTEGER NOT NULL UNIQUE, 
		"session_id"	TEXT NOT NULL,
		"start_date"	DATETIME NOT NULL,
		"expire_date"	DATETIME NOT NULL,
		FOREIGN KEY("user_id") REFERENCES "people"("user_id")
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreatePostsTable()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "posts" (
		"post_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"title"	 VARCHAR(125) NOT NULL UNIQUE,
		"content"	VARCHAR(63,206) NOT NULL,
		"time_creation"	DATETIME NOT NULL,
		PRIMARY KEY("post_id" AUTOINCREMENT),
		FOREIGN KEY("user_id") REFERENCES "people"("user_id")
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreateCommentsTable()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "comments" (
		"comment_id"	INTEGER NOT NULL,
		"user_id"	INTEGER NOT NULL,
		"post_id"	INTEGER NOT NULL,
		"content"	VARCHAR(1000) NOT NULL,
		"time_creation"	DATETIME NOT NULL,
		FOREIGN KEY("user_id") REFERENCES "people"("user_id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		PRIMARY KEY("comment_id")
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreateCategoryTable()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "categories" (
		"category_id"	INTEGER NOT NULL UNIQUE,
		"title"	VARCHAR(125) NOT NULL UNIQUE,
		PRIMARY KEY("category_id" AUTOINCREMENT)
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreatePostCategory()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "post_category" (
		"pc_id"	INTEGER NOT NULL,
		"post_id"	INTEGER NOT NULL,
		"category_id"	INTEGER NOT NULL,
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		FOREIGN KEY("category_id") REFERENCES "categories"("category_id"),
		PRIMARY KEY("pc_id" AUTOINCREMENT)
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreatePostReaction()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "post_reaction" (
		"pr_id"	INTEGER,
		"user_id"	INTEGER NOT NULL,
		"post_id"	INTEGER NOT NULL,
		"reaction"	INTEGER,
		FOREIGN KEY("user_id") REFERENCES "people"("user_id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
		PRIMARY KEY("pr_id")
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
}

func (c *Database) CreateCommentReaction()  {
	stmt, err := c.SqlDb.Prepare(`CREATE TABLE IF NOT EXISTS "comment_reaction" (
		"cr_id"	INTEGER,
		"user_id"	INTEGER NOT NULL,
		"comment_id"	INTEGER NOT NULL,
		"reaction"	INTEGER,
		FOREIGN KEY("comment_id") REFERENCES "comments"("comment_id"),
		FOREIGN KEY("user_id") REFERENCES "people"("user_id"),
		PRIMARY KEY("cr_id")
	);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
 	
}
