package mock

import (
	"time"
	"github.com/cwyang/letsgo/pkg/models"
)

var mockUser = &models.User{
	ID: 1,
	Name: "Kim",
	Email: "kim@mail.com",
	Created: time.Now(),
	Active: true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dup@mail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "kim@mail.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (m *UserModel) ChangePassword(id int, oldpass, newpass string) error {
	return nil
}
