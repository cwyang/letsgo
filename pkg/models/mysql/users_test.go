package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/cwyang/letsgo/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   *models.User
		error  error
	}{
		{
			name: "valid id",
			userID: 1,
			want: &models.User{
				ID: 1,
				Name: "Alice Kim",
				Email: "alice@mail.com",
				Created: time.Date(2023, 10, 05, 18, 31, 22, 0, time.UTC),
				Active: true,
			},
			error: nil,
		},
		{
			name: "zero id",
			userID: 0,
			want: nil,
			error: models.ErrNoRecord,
		},
		{
			name: "notexist id",
			userID: 2,
			want: nil,
			error: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}
			user, err := m.Get(tt.userID)

			if err != tt.error {
				t.Errorf("want %d; got %d", tt.error, err)
			}
			if !reflect.DeepEqual(user, tt.want) {
				t.Errorf("want %v; got %v", tt.want, user)
			}
		})
	}

	
}
