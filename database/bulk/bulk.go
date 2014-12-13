package bulk

import (
	"database/sql"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/jawr/dns/database/connection"
	"strings"
)

// perhaps change to bytes for a binary COPY TO
type Stmt interface {
	Exec(...interface{}) (sql.Result, error)
}

type Item interface {
	BulkInsert(Stmt) error
}

type Insert struct {
	items        []Item
	name         string
	table        string
	args         []string
	sequenceName string
	sequence     int32
	tx           *sql.Tx
	stmt         *sql.Stmt
}

func NewInsert(name, table string, args ...string) (Insert, error) {
	conn, err := connection.Get()
	if err != nil {
		return Insert{}, err
	}
	bi := Insert{}
	serial := uniuri.NewLen(5)
	bi.name = fmt.Sprintf(name, serial)
	bi.table = fmt.Sprintf(table, bi.name)
	tx, err := conn.Begin()
	if err != nil {
		return bi, err
	}
	// close up our defer tx.Rollback()
	bi.tx = tx
	_, err = bi.tx.Exec(bi.table)
	if err != nil {
		return bi, err
	}
	// prepare args
	for idx, i := range args {
		args[idx] = fmt.Sprintf(`"%s"`, i)
	}
	bi.args = args
	stmt := fmt.Sprintf(
		`COPY %s (%s) FROM STDIN`,
		bi.name,
		strings.Join(bi.args, ", "),
	)
	bi.stmt, err = bi.tx.Prepare(stmt)
	if err != nil {
		return bi, err
	}
	return bi, nil
}

func (bi *Insert) Add(i Item) error {
	return i.BulkInsert(bi.stmt)
}

func (bi *Insert) Insert() error {
	_, err := bi.stmt.Exec()
	if err != nil {
		return err
	}
	err = bi.stmt.Close()
	return err
}

func (bi *Insert) Merge(query string) error {
	query = fmt.Sprintf(query, bi.name)
	_, err := bi.tx.Exec(query)
	return err
}

func (bi *Insert) Index(query string) error {
	query = fmt.Sprintf(query, bi.name)
	_, err := bi.tx.Exec(query)
	return err
}

func (bi *Insert) Finish() error {
	return bi.tx.Commit()
}

func (bi *Insert) Close() error {
	// does this matter as it's called in a defer?
	bi.stmt.Close()
	return bi.tx.Rollback()
}
