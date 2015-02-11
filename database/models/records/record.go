package records

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records/types"
	"github.com/jawr/dns/database/models/tlds"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"strconv"
	"strings"
	"time"
)

type Args struct {
	TTL  uint     `json:"ttl"`
	Args []string `json:"args,omitempty,omitempty"`
}

type Record struct {
	UUID   uuid.UUID      `json:"uuid"`
	Domain domains.Domain `json:"domain"`
	Name   string         `json:"name"`
	Args   Args           `json:"args"`
	Type   types.Type     `json:"type"`
	Date   time.Time      `json:"parse_date"`
	Added  time.Time      `json:"added"`
}

func New(line, origin string, tld tlds.TLD, ttl uint, date time.Time) (Record, error) {
	fields := strings.Fields(line)
	r := Record{}

	name := fields[0]
	if !strings.HasSuffix(name, ".") {
		name += "." + origin
	}
	r.Domain = domains.New(name, tld)
	r.Name = name
	if name == domains.CleanDomain(r.Name, tld) {
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
	rt, err := types.Detect(fields[typeIdx])
	if err != nil {
		return r, errors.New("Unable to detect Record Type.")
	}
	r.Type, err = types.New(rt, tld)
	r.Args.Args = fields[typeIdx+1:]
	r.Date = date
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(
		fmt.Sprintf("%v", r),
	))
	r.UUID = id
	r.Added = time.Now()
	return r, nil
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
			r.Type.ID,
			r.Domain.TLD.ID,
		),
		r.UUID.String(),
		r.Domain.UUID.String(),
		r.Name,
		string(args),
		r.Type.ID,
		r.Date,
	)
	return err
}
