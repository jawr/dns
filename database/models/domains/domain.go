package domains

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tlds"
	"strings"
)

type Domain struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	TLD  tlds.TLD  `json:"tld"`
}

func New(name string, t tlds.TLD) Domain {
	name = CleanDomain(name, t)
	args := strings.Split(name, ".")
	name = args[len(args)-1]
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(fmt.Sprintf("%s_%d", name, t.ID)))
	return Domain{
		UUID: id,
		Name: name,
		TLD:  t,
	}
}

func (d Domain) UID() string {
	return d.UUID.String()
}

func (d Domain) String() string {
	return d.Name + "." + d.TLD.Name
}

func (d Domain) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO domain (uuid, name, tlds) VALUES ($1, $2, $3)",
		d.UUID.String(),
		d.Name,
		d.TLD.ID,
	)
	return err
}

func CleanDomain(s string, t tlds.TLD) string {
	s = strings.TrimSuffix(s, ".")
	s = strings.TrimSuffix(s, "."+t.Name)
	args := strings.Split(s, ".")
	return args[len(args)-1]
}

func (d Domain) CleanSubdomain(s string) string {
	s = strings.TrimSuffix(s, ".")
	s = strings.TrimSuffix(s, "."+d.TLD.Name)
	s = strings.TrimSuffix(s, "."+d.Name)
	return s
}
