package connection

import (
	"github.com/jackc/pgx"
	"log"
)

type connection struct {
	db *pgx.ConnPool
}

type Row interface {
	Scan(dest ...interface{}) error
}

var single *connection = nil

func Get() (*pgx.Conn, error) {
	pool, err := pgx.NewConnPool(extractConfig())
	if err != nil {
		log.Printf("Unable to connect to database: %s", err)
		return nil, err
	}
	return pool.Acquire()
}

func extractConfig() pgx.ConnPoolConfig {
	var config pgx.ConnPoolConfig
	config.Host = "localhost"
	config.User = "dns"
	config.Password = "dns!pass$"
	config.Database = "dns"
	return config
}
