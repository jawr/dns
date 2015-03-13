package users

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
)

const (
	SELECT string = "SELECT * FROM users "
)

type Result struct {
	One  func() (User, error)
	List func() ([]User, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (User, error) {
			return Get(query, args...)
		},
		func() ([]User, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByEmail(email string) Result {
	return newResult(SELECT+"WHERE email = $1", email)
}

func GetByID(id int) Result {
	return newResult(SELECT+"WHERE id = $1", id)
}

func parseRow(row connection.Row) (User, error) {
	u := User{}
	var settings []byte
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Added, &u.Updated, &settings)
	if err != nil {
		return u, err
	}
	err = json.Unmarshal(settings, &u.Settings)
	return u, err
}

func Get(query string, args ...interface{}) (User, error) {
	conn, err := connection.Get()
	if err != nil {
		return User{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]User, error) {
	conn, err := connection.Get()
	if err != nil {
		return []User{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []User{}, err
	}
	var list []User
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}
