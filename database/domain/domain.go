package domain

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jawr/dns/database/bulk"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/tld"
	"strings"
)

type Domain struct {
	UUID uuid.UUID
	Name string
	TLD  tld.TLD
}

func New(name string, t tld.TLD) Domain {
	name = strings.TrimSuffix(name, ".")
	name = strings.TrimSuffix(name, "."+t.Name)
	uuid := uuid.NewSHA1(uuid.NameSpace_OID, []byte(fmt.Sprintf("%s_%d", name, t.ID)))
	return Domain{
		UUID: uuid,
		Name: name,
		TLD:  t,
	}
}

func (d *Domain) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO domain (uuid, name, tld) VALUES ($1, $2, $3)",
		d.UUID.String(),
		d.Name,
		d.TLD.ID,
	)
	return err
}

func NewBulkInsert() (bulk.Insert, error) {
	table := `CREATE TEMP TABLE %s (
		uuid UUID,
		name TEXT,
		tld INT
	) ON COMMIT DROP`
	tableName := "_domain__%s"
	bi, err := bulk.NewInsert(tableName, table, "uuid", "name", "tld")
	return bi, err
}

func (d *Domain) BulkInsert(stmt bulk.Stmt) error {
	_, err := stmt.Exec(d.UUID.String(), d.Name, d.TLD.ID)
	return err
}

func GetByNameAndTLD() string {
	return "SELECT * FROM domain WHERE name = $1 AND tld = $2"
}

func GetByID() string {
	return "SELECT * FROM domain WHERE id = $1"
}

func parseRow(row connection.Row) (Domain, error) {
	d := Domain{}
	var tldId int32
	err := row.Scan(&d.UUID, &d.Name, &tldId)
	if err != nil {
		return d, err
	}
	t, err := tld.Get(tld.GetByID(), tldId)
	d.TLD = t
	return d, nil
}

func Get(query string, args ...interface{}) (Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return Domain{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Domain, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Domain{}, err
	}
	rows, err := conn.Query(query)
	defer rows.Close()
	if err != nil {
		return []Domain{}, err
	}
	var list []Domain
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}
