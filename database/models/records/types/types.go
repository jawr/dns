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
	if _, ok := byNameAndTLD[name]; ok {
		if t, ok := byNameAndTLD[name][tld.ID]; ok {
			return t, nil
		}
	}
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
	if err == nil {
		if _, ok := byNameAndTLD[name]; !ok {
			byNameAndTLD[name] = make(map[int32]Type)
		}
		byNameAndTLD[name][tld.ID] = rt
	}
	return rt, err
}
