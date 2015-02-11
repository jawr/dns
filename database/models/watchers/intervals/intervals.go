package intervals

import (
	"github.com/jawr/dns/database/connection"
)

type Interval struct {
	ID    int32  `json:"id"`
	Value string `json:"value"`
}

func New(value string) (Interval, error) {
	i := Interval{}
	i.Value = value
	conn, err := connection.Get()
	if err != nil {
		return i, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_interval($1)", value).Scan(&id)
	i.ID = id
	return i, err
}
