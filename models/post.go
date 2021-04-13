package models

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type Post struct {
	Header string `db:"header" json:"header"`
	Body   string `db:"body" json:"body"`
}

const (
	ErrorMissingHeader = "missing TODO header test"
	ErrorMissingBody   = "missing TODO body test"
	InsertStatement    = "INSERT INTO todo_post(header, body) VALUES ($1, $2)"
	SelectAllStatement = "SELECT header, body FROM todo_post"
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

		err = rows.Scan(&post.Header, &post.Body)
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
