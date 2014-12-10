package domain

import (
	"github.com/jawr/dns/database/bulk"
)

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
