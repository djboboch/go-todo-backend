package models

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type Post struct {
	Id     string `json:"id"`
	Header string `db:"header" json:"header"`
	Body   string `db:"body" json:"body"`
}

const (
	ErrorMissingHeader = "missing TODO header test"
	ErrorMissingBody   = "missing TODO body test"
	InsertStatement    = "INSERT INTO todo_post(header, body) VALUES ($1, $2)"
	SelectAllStatement = "SELECT id, header, body FROM todo_post"
	DeleteStatement    = "DELETE FROM todo_post WHERE id = $1"
)

func CreatePost(db *sql.DB, p Post) error {
	var err error
	if p.Header == "" {
		return errors.New(ErrorMissingHeader)
	}

	if p.Body == "" {
		return errors.New(ErrorMissingBody)
	}

	_, err = db.Exec(InsertStatement, p.Header, p.Body)
	if err != nil {
		return err
	}

	return nil
}

func AllPosts(db *sql.DB) ([]Post, error) {
	var err error

	rows, err := db.Query(SelectAllStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.Id, &post.Header, &post.Body)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

//DeletePost deletes the post with the passed in ID
func DeletePost(db *sql.DB, id string) error {
	var err error

	_, err = db.Exec(DeleteStatement, id)

	if err != nil {
		return err
	}

	return nil
}
