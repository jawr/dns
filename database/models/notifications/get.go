package notifications

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/users"
)

const (
	SELECT string = "SELECT * FROM notification "
)

type Result struct {
	One  func() (Notification, error)
	List func() ([]Notification, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Notification, error) {
			return Get(query, args...)
		},
		func() ([]Notification, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByUser(user users.User) Result {
	return newResult(SELECT+"WHERE user_id = $1", user.ID)
}

func Get(query string, args ...interface{}) (Notification, error) {
	conn, err := connection.Get()
	if err != nil {
		return Notification{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Notification, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Notification{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Notification{}, err
	}
	var list []Notification
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}

func parseRow(row connection.Row) (Notification, error) {
	n := Notification{}
	var userID int
	var messages, archived []byte
	err := row.Scan(&userID, &messages, &n.Updated, &n.Alerts, &archived)
	if err != nil {
		return n, err
	}
	err = json.Unmarshal(archived, &n.Archived)
	if err != nil {
		return n, err
	}
	err = json.Unmarshal(messages, &n.Messages)
	if err != nil {
		return n, err
	}
	n.User, err = users.GetByID(userID).One()
	if err != nil {
	}
	return n, err
}
