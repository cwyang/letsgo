package mysql

import (
	"database/sql"
	"errors"
	
	"github.com/cwyang/letsgo/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(title, content, expires string) (int, error) {
	stmt := `insert into notes (title, content, created, expires)
		values(?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *NoteModel) Get(id int) (*models.Note, error) {
	stmt := `select id, title, content, created, expires from notes
                 where expires > utc_timestamp() and id = ?`
	row := m.DB.QueryRow(stmt, id)

	n := &models.Note{}

	err := row.Scan(&n.ID, &n.Title, &n.Content, &n.Created, &n.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return n, nil
}

// 10 most recently creates notes
func (m *NoteModel) Latest() ([]*models.Note, error) {
	stmt := `select id, title, content, created, expires from notes
                 where expires > utc_timestamp() order by created desc limit 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*models.Note{} // slice of object pointers

	for rows.Next() {
		n := &models.Note{}
		err = rows.Scan(&n.ID, &n.Title, &n.Content, &n.Created, &n.Expires)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
