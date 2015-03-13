// Package connection is a singleton object used to connect to a Postgres instance. The package
// relies on a config.json file to read it's globals from.
package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stathat/jconfig"
)

// Row is an interface that allows us to pass sql.Rows and sql.Row to our model specific
// parseRow functions.
type Row interface {
	Scan(dest ...interface{}) error
}

// Our singleton.
type connection struct {
	db *sql.DB
}

var single *connection = nil

// Get attempts to setup our connection, if already setup it checks if still alive and
// then attempts to retrieve a connection.
//
// Warning, this method can panic (call to setup).
func Get() (*sql.DB, error) {
	if single == nil {
		single = new(connection)
		if err := single.setup(); err != nil {
			return single.db, err
		}
	}
	err := single.db.Ping()
	if err != nil {
		if err := single.setup(); err != nil {
			return single.db, err
		}
	}
	return single.db, nil
}

// setup uses config.json to try and connect to Postgres. Warning, this method can panic.
func (c *connection) setup() error {
	config := jconfig.LoadConfig("config.json")
	user := config.GetString("db_user")
	pass := config.GetString("db_pass")
	name := config.GetString("db_name")
	if len(user+pass+name) == 0 {
		panic("error setup config file for database connection")
	}
	conn, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			user,
			pass,
			name,
		),
	)
	if err != nil {
		return err
	}
	c.db = conn
	return nil
}
