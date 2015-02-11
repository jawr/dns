package types

import (
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tlds"
)

type Type struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func New(name string, tld tlds.TLD) (Type, error) {
	conn, err := connection.Get()
	if err != nil {
		return Type{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT ensure_record_table($1, $2)", name, tld.ID).Scan(&id)
	rt := Type{
		ID:   id,
		Name: name,
	}
	return rt, err
}

func (rt Type) UID() string { return rt.Name }
