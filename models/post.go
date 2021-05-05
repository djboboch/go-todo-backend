package models

import (
	"database/sql"
	"errors"
	"github.com/segmentio/ksuid"

	_ "github.com/lib/pq"
)

type Post struct {
	Id             string `json:"id"`
	Content        string `db:"content" json:"content"`
	IsItemFinished bool   `db:"isItemFinished" json:"isItemFinished"`
}

const (
	ErrorMissingContent = "missing TODO content text"
	InsertStatement     = "INSERT INTO post_item(id, content) VALUES ($1, $2);"
	SelectAllStatement  = "SELECT id, content, isItemFinished FROM post_item;"
	UpdateStatement     = "UPDATE post_item SET content = $2, isItemFinished = $3 WHERE id = $1;"
	DeleteStatement     = "DELETE FROM post_item WHERE id = $1;"
)

func CreatePost(db *sql.DB, content string) error {
	var err error

	if content == "" {
		return errors.New(ErrorMissingContent)
	}

	_, err = db.Exec(InsertStatement, ksuid.New().String(), content)
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

		err = rows.Scan(&post.Id, &post.Content, &post.IsItemFinished)
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

func UpdatePost(db *sql.DB, post Post) error {
	var err error

	_, err = db.Exec(UpdateStatement, post.Id, post.Content, post.IsItemFinished)

	if err != nil {
		return err
	}

	return nil
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
