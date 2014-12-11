package record

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/domain"
	"github.com/jawr/dns/database/record_type"
	"time"
)

type RecordArgs struct {
	TTL  uint     `json:"ttl"`
	Args []string `json:"args,omitempty,omitempty"`
}

type Record struct {
	UUID       uuid.UUID
	Domain     domain.Domain
	Name       string
	Args       RecordArgs
	RecordType record_type.RecordType
	Date       time.Time
	Added      time.Time
}

func New(name string, date time.Time, d domain.Domain, args RecordArgs, rt record_type.RecordType) Record {
	name = d.CleanSubdomain(name)
	r := Record{
		Domain:     d,
		Name:       name,
		Args:       args,
		RecordType: rt,
		Date:       date,
	}
	// probably a nicer way to do this
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(
		fmt.Sprintf("%v", r),
	))
	r.UUID = id
	// take from parser filename rather than assume parse time
	r.Added = time.Now()
	return r
}

func (r Record) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	args, err := json.Marshal(r.Args)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`INSERT INTO record
		(uuid, domain, name, args, record_type, parser_date) VALUES
		($1,   $2,     $3,   $4,   $5,          $6)`,
		r.UUID.String(),
		r.Domain.UUID.String(),
		r.Name,
		string(args),
		r.RecordType.ID,
		r.Date,
	)
	return err
}
