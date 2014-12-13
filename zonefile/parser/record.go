package parser

import (
	"errors"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/database/models/record_type"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/util"
	"log"
	"strconv"
	"strings"
)

type Record struct {
	Name       string
	Args       []string
	TTL        uint
	RecordType string
	TLD        tld.TLD
}

// refactor this "package"

func (p Parser) buildRecordRow(r Record) (record.Record, error) {
	args := record.RecordArgs{
		TTL:  r.TTL,
		Args: r.Args,
	}

	rt, err := record_type.New(r.RecordType, p.tld)
	if err != nil {
		log.Printf("ERROR: Record.Save:record_type.New: %s", err)
		return record.Record{}, err
	}

	d := domain.New(r.Name, r.TLD)
	err = p.domainInsert.Add(&d)
	if err != nil {
		log.Printf("ERROR: Record.Save:domainInsert.Add: %s", err)
		return record.Record{}, err
	}

	rr := record.New(r.Name, p.date, d, args, rt)

	if rr.Name == domain.CleanDomain(r.Name, d.TLD) {
		rr.Name = "@"
	}
	return rr, nil
}

// This function assumes the following:
// <name> <ttl?> <type> <args>
func (p Parser) getRecord(fields []string) (record.Record, error) {
	//defer util.Un(util.Trace())
	r := Record{}
	r.Name = fields[0]
	if !strings.HasSuffix(r.Name, ".") {
		r.Name += "." + p.origin
	}
	r.TTL = p.ttl
	typeIdx := 1
	// strip \sin\s
	if len(fields) > 3 {
		fields = util.FilterIN(fields)
	}
	if len(fields) > 3 {
		ttl, err := strconv.ParseUint(fields[1], 10, 0)
		if err == nil {
			typeIdx = 2
			r.TTL = uint(ttl)
		} else {
			// detect and fix these, maybe go to their own channel/table
			log.Printf("WARN: getRecord:parseTTL: len(fields) == %d, fields: %s", len(fields), fields)
		}
	}
	if len(fields) <= typeIdx {
		log.Printf("ERROR: getRecord:setTypeIdx: len(fields) == %d, fields: %s", len(fields), fields)
		return record.Record{}, errors.New("Unable to set typeIdx in getRecord.")
	}
	var err error
	r.RecordType, err = detectRecordType(fields[typeIdx])
	if err != nil {
		log.Printf("ERROR: getRecord.parseType: len(fields) == %d, fields: %s", len(fields), fields)
		return record.Record{}, errors.New("Unable to detect Record Type.")
	}
	r.Args = fields[typeIdx+1:]
	r.TLD = p.tld
	return p.buildRecordRow(r)
}
