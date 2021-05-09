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

type PostModel struct {
	DB *sql.DB
}

const (
	ErrorMissingContent = "missing TODO content text"
	InsertStatement     = "INSERT INTO post_item(id, content) VALUES ($1, $2);"
	SelectAllStatement  = "SELECT id, content, isItemFinished FROM post_item;"
	UpdateStatement     = "UPDATE post_item SET content = $2, isItemFinished = $3 WHERE id = $1;"
	DeleteStatement     = "DELETE FROM post_item WHERE id = $1;"
)

func (m PostModel) Create(content string) (*Post, error) {
	var err error

	if content == "" {
		return nil, errors.New(ErrorMissingContent)
	}

	post := Post{
		Id:             ksuid.New().String(),
		Content:        content,
		IsItemFinished: false,
	}

	_, err = m.DB.Exec(InsertStatement, post.Id, post.Content)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (m PostModel) All() ([]Post, error) {
	var err error

	rows, err := m.DB.Query(SelectAllStatement)
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

func (m PostModel) Update(post Post) error {
	var err error

	_, err = m.DB.Exec(UpdateStatement, post.Id, post.Content, post.IsItemFinished)

	if err != nil {
		return err
	}

	return nil
}

//Delete deletes the post with the passed in ID
func (m PostModel) Delete(id string) error {
	var err error

	_, err = m.DB.Exec(DeleteStatement, id)

	if err != nil {
		return err
	}

	return nil
}
