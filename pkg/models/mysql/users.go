package mysql

import (
	"database/sql"
	"errors"
	"strings"
	"fmt"
	
	"github.com/cwyang/letsgo/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `insert into users (name, email, hashed_password, created)
                 values(?, ?, ?, utc_timestamp())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 &&
				strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "select id, hashed_password from users where email = ? and active = true"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	
	return id, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `select id, name, email, hashed_password, created, active from users where id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name,
		&u.Email, &u.HashedPassword, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
func (m *UserModel) ChangePassword(id int, oldpass, newpass string) error {
	u, err := m.Get(id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(oldpass))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		} else {
			return err
		}
	}
	fmt.Printf("2\n")
	u.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(newpass), 12)
	if err != nil {
		return err
	}
	fmt.Printf("3\n")
	
	stmt := `update users set hashed_password = ? where id = ?`
	_, err = m.DB.Exec(stmt, string(u.HashedPassword), u.ID)
	return err
}
