package tlds

import (
	"github.com/jawr/dns/database/connection"
	"strings"
)

type TLD struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func New(name string) (TLD, error) {
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_tld($1)", name).Scan(&id)
	tld := TLD{
		ID:   id,
		Name: name,
	}
	return tld, err
}

func (t TLD) UID() string { return t.Name }

func Detect(s string) (TLD, error) {
	tld, err := GetByName(s).One()
	if err != nil {
		args := strings.Split(s, ".")
		if len(args) > 1 {
			return Detect(strings.Join(args[1:], "."))
		}
		return tld, err
	}
	return tld, err
}

func DetectDomainAndTLD(s string) (string, TLD, error) {
	tld, err := Detect(s)
	s = strings.TrimSuffix(s, "."+tld.Name)
	args := strings.Split(s, ".")
	return args[len(args)-1], tld, err
}
