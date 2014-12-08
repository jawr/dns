package bulk

import (
	"database/sql"
	"fmt"
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
	sequenceName string
	sequence     int32
	tx           *sql.Tx
	stmt         *sql.Stmt
}

var serial int32 = 0

func NewInsert(name, table, sequenceName string, args ...string) (Insert, error) {
	conn, err := connection.Get()
	if err != nil {
		return Insert{}, err
	}
	bi := Insert{}
	serial++
	bi.name = fmt.Sprintf(name, serial)
	bi.table = fmt.Sprintf(table, bi.name)
	bi.sequenceName = sequenceName
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
	err = bi.tx.QueryRow("SELECT nextval($1)", bi.sequenceName).Scan(&bi.sequence)
	if err != nil {
		return bi, err
	}
	// prepare args
	for idx, i := range args {
		args[idx] = fmt.Sprintf(`"%s"`, i)
	}
	stmt := fmt.Sprintf(
		`COPY "%s" (%s) FROM STDIN`,
		bi.name,
		strings.Join(args, ", "),
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

func (bi *Insert) Close() error {
	return bi.tx.Commit()
}
