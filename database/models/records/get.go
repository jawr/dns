package records

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records/types"
	"github.com/jawr/dns/database/models/zonefile/parser"
)

const (
	SELECT string = "SELECT * FROM record "
)

type Result struct {
	One  func() (Record, error)
	List func() ([]Record, error)
}

func newResult(query string, args ...interface{}) Result {
	return Result{
		func() (Record, error) {
			return Get(query, args...)
		},
		func() ([]Record, error) {
			return GetList(query, args...)
		},
	}
}

func GetAll() Result {
	return newResult(SELECT)
}

func GetByUUID(uuid string) Result {
	return newResult(SELECT+"WHERE uuid = $1", uuid)
}

func GetByDomain(domain domains.Domain) Result {
	return newResult(SELECT+"WHERE domain = $1", domain.UUID.String())
}

func parseRow(row connection.Row) (Record, error) {
	r := Record{}
	var rUUID, dUUID string
	var rtID int32
	var pID int32
	var jsonBuf []byte
	err := row.Scan(&rUUID, &dUUID, &r.Name, &jsonBuf, &rtID, &r.Date, &pID, &r.Added)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(jsonBuf, &r.Args)
	if err != nil {
		return r, err
	}
	r.UUID = uuid.Parse(rUUID)
	r.Domain, err = domains.GetByUUID(dUUID).One()
	if err != nil {
		return r, err
	}
	r.Type, err = types.GetByID(rtID).One()
	if err != nil {
		return r, err
	}
	// ignore error as we dont need a Parser
	p, _ := parser.Get(parser.GetByID(), pID)
	r.Parser = p
	return r, nil
}

func Get(query string, args ...interface{}) (Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return Record{}, err
	}
	row := conn.QueryRow(query, args...)
	return parseRow(row)
}

func GetList(query string, args ...interface{}) ([]Record, error) {
	conn, err := connection.Get()
	if err != nil {
		return []Record{}, err
	}
	rows, err := conn.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return []Record{}, err
	}
	var list []Record
	for rows.Next() {
		rt, err := parseRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, rt)
	}
	return list, rows.Err()
}
