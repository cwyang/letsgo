package mock

import (
	"time"
	"github.com/cwyang/letsgo/pkg/models"
)

var mockNote = &models.Note{
	ID: 1,
	Title: "Mock Note",
	Content: "Quick Brown Fox Jumps Over the Little Lazy Dog",
	Created: time.Now(),
	Expires: time.Now(),
}

type NoteModel struct{}

func (m *NoteModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *NoteModel) Get(id int) (*models.Note, error) {
	switch id {
	case 1:
		return mockNote, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *NoteModel) Latest() ([]*models.Note, error) {
	return []*models.Note{mockNote}, nil
}
