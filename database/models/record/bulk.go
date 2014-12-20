package record

import (
	"encoding/json"
	"github.com/jawr/dns/database/bulk"
)

func NewBulkInsert() (bulk.Insert, error) {
	table := `CREATE TEMP TABLE %s (
		uuid UUID,
		domain UUID,
		name TEXT,
		args jsonb,
		record_type INT,
		parser_date DATE
	) ON COMMIT DROP`
	tableName := "_record__%s"
	bi, err := bulk.NewInsert(tableName, table, "uuid", "domain", "name", "args", "record_type", "parser_date")
	return bi, err
}

func (r *Record) BulkInsert(stmt bulk.Stmt) error {
	args, err := json.Marshal(r.Args)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.UUID.String(), r.Domain.UUID.String(), r.Name, string(args), r.RecordType.ID, r.Date)
	return err
}
