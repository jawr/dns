package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stathat/jconfig"
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
		if err := single.setup(); err != nil {
			return single, err
		}
	}
	err := single.db.Ping()
	return single.db, err
}

func (c *connection) setup() error {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
			user,
			pass,
			name,
			host,
		),
	)
	if err != nil {
		return err
	}
	c.db = conn
	return nil
}

var host, user, pass, name string

func init() {
	config := jconfig.LoadConfig("config.json")
	host = config.GetString("db_host")
	user = config.GetString("db_user")
	pass = config.GetString("db_pass")
	name = config.GetString("db_name")
	if len(host+user+pass+name) == 0 {
		panic("error setup config file for database connection")
	}
}
