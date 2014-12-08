package connection

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type connection struct {
	db *sql.DB
}

type Row interface {
	Scan(dest ...interface{}) error
}

var single *connection = nil

func Get() (*sql.DB, error) {
	if single == nil {
		single = new(connection)
		single.setup()
	}
	err := single.db.Ping()
	return single.db, err
}

func (c *connection) setup() error {
	conn, err := sql.Open("postgres", "user=dns password=dns!pass$ dbname=dns")
	if err != nil {
		log.Printf("ERROR: connection.go:setup:sql.Open: %s", err)
		return err
	}
	c.db = conn
	return nil
}
