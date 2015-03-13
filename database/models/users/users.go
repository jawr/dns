package users

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/jawr/dns/database/connection"
	"time"
)

type Settings struct {
	IsAdmin bool `json:"is_admin,omitempty"`
}

type User struct {
	ID       int32     `json:"id"`
	Email    string    `json:"email"`
	Password []byte    `json:"password,omit"`
	Added    time.Time `json:"added"`
	Updated  time.Time `json:"updated"`
	Settings Settings  `json:"settings:`
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func PasswordCrypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func New(email, password string) (User, error) {
	if u, err := GetByEmail(email).One(); err == nil {
		if CheckPassword(email, password) {
			return u, nil
		}
		return User{}, errors.New("Email does not exists, or password does not match.")
	}
	u := User{}
	pass, err := PasswordCrypt([]byte(password))
	if err != nil {
		return u, err
	}
	u.Email = email
	u.Password = pass
	conn, err := connection.Get()
	if err != nil {
		return u, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_user($1, $2)", email, pass).Scan(&id)
	if err == nil {
		u.ID = id
	}
	return u, err
}

func CheckPassword(email, password string) bool {
	conn, err := connection.Get()
	if err != nil {
		return false
	}
	var pass []byte
	err = conn.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&pass)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return false
	}
	return true
}
