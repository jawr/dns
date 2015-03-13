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
	"github.com/jawr/dns/database/models/zonefile/parser"
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
	Parser parser.Parser  `json:"parser"`
	Added  time.Time      `json:"added"`
}

// New creates a new Record, it takes the entire resource record line, the origin of the zonefile, the tld associated, the ttl and a parse date. It then creates a formalized Record.
func New(line, origin string, tld tlds.TLD, ttl uint, date time.Time, parserID int32) (Record, error) {
	fields := strings.Fields(line)
	r := Record{}

	name := fields[0]
	// set origin if it's not already there
	if !strings.HasSuffix(name, ".") {
		name += "." + origin
	}
	// strip domain name
	name = domains.CleanDomain(name, tld)
	r.Domain = domains.New(name, tld)
	name = strings.TrimSuffix(name, ".")
	name = strings.TrimSuffix(name, r.Domain.TLD.Name)
	name = strings.TrimSuffix(name, r.Domain.Name)
	name = strings.TrimSuffix(name, ".")
	// check if we are referencing top level
	if len(name) == 0 {
		name = "@"
	}
	r.Name = name

	typeIdx := 1
	if len(fields) > 3 {
		fields = util.FilterIN(fields)
	}
	if len(fields) > 3 {
		ttlFromFields, err := strconv.ParseUint(fields[1], 10, 64)
		if err == nil {
			typeIdx = 2
			ttl = uint(ttlFromFields)
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
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(
		fmt.Sprintf("%v%d", r, parserID),
	))
	r.Date = date
	r.Args.TTL = ttl
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
				(uuid, domain, name, args, record_type, parser_date, parser) VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
			r.Type.ID,
			r.Domain.TLD.ID,
		),
		r.UUID.String(),
		r.Domain.UUID.String(),
		r.Name,
		string(args),
		r.Type.ID,
		r.Date,
		r.Parser.ID,
	)
	return err
}
