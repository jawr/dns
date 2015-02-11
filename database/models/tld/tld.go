package tlds

import (
	"fmt"
	"github.com/jawr/dns/database/cache"
	"github.com/jawr/dns/database/connection"
	"net/url"
	"reflect"
	"strings"
)

type TLD struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

var c = cache.New()
var cacheGetID = cache.NewCacheInt32()

func New(name string) (TLD, error) {
	if t, ok := c.Check(name); ok {
		return t.(TLD), nil
	}
	conn, err := connection.Get()
	if err != nil {
		return TLD{}, err
	}
	var id int32
	err = conn.QueryRow("SELECT insert_tld($1)", name).Scan(&id)
	t := TLD{
		ID:   id,
		Name: name,
	}
	c.Add(t)
	return t, err
}

func (t TLD) UID() string { return t.Name }

func Detect(s string) (TLD, error) {
	t, err := Get(GetByName(), s)
	if err != nil {
		args := strings.Split(s, ".")
		if len(args) > 1 {
			return Detect(strings.Join(args[1:], "."))
		}
		return TLD{}, err
	}
	return t, err
}

func DetectDomainAndTLD(s string) (string, TLD, error) {
	t, err := Detect(s)
	s = strings.TrimSuffix(s, "."+t.Name)
	args := strings.Split(s, ".")
	return args[len(args)-1], t, err
}
