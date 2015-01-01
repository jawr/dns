package record

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record_type"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"strconv"
	"strings"
	"time"
)

type RecordArgs struct {
	TTL  uint     `json:"ttl"`
	Args []string `json:"args,omitempty,omitempty"`
}

type Record struct {
	UUID       uuid.UUID              `json:"uuid"`
	Domain     domain.Domain          `json:"domain"`
	Name       string                 `json:"name"`
	Args       RecordArgs             `json:"args"`
	RecordType record_type.RecordType `json:"type"`
	Date       time.Time              `json:"parse_date"`
	Added      time.Time              `json:"added"`
}

func New(line, origin string, t tld.TLD, ttl uint, date time.Time) (Record, error) {
	fields := strings.Fields(line)
	r := Record{}

	name := fields[0]
	if !strings.HasSuffix(name, ".") {
		name += "." + origin
	}
	r.Domain = domain.New(name, t)
	r.Name = name
	if name == domain.CleanDomain(r.Name, t) {
		r.Name = "@"
	}
	r.Args.TTL = ttl

	typeIdx := 1
	if len(fields) > 3 {
		fields = util.FilterIN(fields)
	}
	if len(fields) > 3 {
		ttl, err := strconv.ParseUint(fields[1], 10, 0)
		if err == nil {
			typeIdx = 2
			r.Args.TTL = uint(ttl)
		} else {
			log.Warn("Unable to parse RR: len(fields) == %d, fields: %s", len(fields), fields)
		}
	}
	if len(fields) <= typeIdx {
		return r, errors.New("Unable to set typeIdx in getRecord.")
	}
	rt, err := record_type.DetectRecordType(fields[typeIdx])
	if err != nil {
		return r, errors.New("Unable to detect Record Type.")
	}
	r.RecordType, err = record_type.New(rt, t)
	r.Args.Args = fields[typeIdx+1:]
	r.Date = date
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(
		fmt.Sprintf("%v", r),
	))
	r.UUID = id
	r.Added = time.Now()
	return r, nil
}

func NewOld(name string, date time.Time, d domain.Domain, args RecordArgs, rt record_type.RecordType) Record {
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
	_, err = conn.Exec(
		fmt.Sprintf(`INSERT INTO record__%d_%d
				(uuid, domain, name, args, record_type, parser_date) VALUES
				($1, $2, $3, $4, $5, $6)`,
			r.RecordType.ID,
			r.Domain.TLD.ID,
		),
		r.UUID.String(),
		r.Domain.UUID.String(),
		r.Name,
		string(args),
		r.RecordType.ID,
		r.Date,
	)
	return err
}

func (r *RecordArgs) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), r)
}

func (r *RecordArgs) Unmarshal(data []byte) error {
	return nil
}
