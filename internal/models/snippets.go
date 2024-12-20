package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	statement := `SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() LIMIT 10`

	result, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var snippets []Snippet

	for result.Next() {
		var s Snippet
		err = result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title,content,created,expires)
	VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`

	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (m *SnippetModel) Get(id int) (Snippet, error) {
	statement := `SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	result := m.DB.QueryRow(statement, id)
	var s Snippet
	err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}

	}
	return s, nil
}
