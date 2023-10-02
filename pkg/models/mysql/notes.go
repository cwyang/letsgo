package mysql

import (
	"database/sql"
	
	"github.com/cwyang/letsgo/pkg/models"
)

type NotesModel struct {
	DB *sql.DB
}

func (m *NotesModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

func (m *NotesModel) Get(id int) (*models.Notes, error) {
	return nil, nil
}

// 10 most recently creates notes
func (m *NotesModel) Latest(id int) (*models.Notes, error) {
	return nil, nil
}
