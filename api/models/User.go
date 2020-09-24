package models

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	ID       uint32 `json:"id"`
	Login    string `json:"login" pg:",unique,notnull"`
	Password string `json:"password" pg:",notnull"`
}

type UserError struct {
	Err error
}

func (r *UserError) Error() string {
	return r.Err.Error()
}

func FindUserByLogin(db *pg.DB, login string) (*User, error) {
	var users []User
	err := db.Model(&users).Where("login = ?", login).First()
	if err != nil {
		log.Printf("Unexpected error: ", err)
		return nil, err
	}

	if len(users) == 0 {
		log.Printf("User '%s' not found ", login)
		return nil, nil
	}
	user := &users[0]
	log.Printf("Found user %d", user.ID)
	return user, nil
}

func CreateUser(db *pg.DB, user *User) error {
	result, err := db.
		Model(user).
		OnConflict("DO NOTHING").
		Insert()

	if err != nil {
		log.Printf("Unexpected error: ", err)
		return err
	}

	if result.RowsAffected() == 0 {
		log.Printf("Failed to create user. Login '%s' already taken", user.Login)
		return &UserError{Err: errors.New("login already used")}
	}
	log.Printf("User %d '%s' created", user.ID, user.Login)
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %s", err)
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
