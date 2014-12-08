package parser

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/jawr/dns/database/domain"
	"github.com/jawr/dns/database/record"
	"github.com/jawr/dns/database/record_type"
	"github.com/jawr/dns/database/tld"
	"io"
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

func (r Record) Save() {
	un(trace())
	var rArgs []string
	if len(r.Args) > 1 {
		rArgs = r.Args[1:]
	}
	args := record.RecordArgs{
		TTL:  r.TTL,
		Addr: r.Args[0],
		Args: rArgs,
	}

	h := md5.New()
	io.WriteString(h, r.Name)
	io.WriteString(h, r.Args[0])
	io.WriteString(h, r.RecordType)
	hash := hex.EncodeToString(h.Sum(nil))

	recordType, err := record_type.New(r.RecordType)
	if err != nil {
		log.Printf("ERROR: Record.Save:record_type.New: %s", err)
		return
	}

	domain, err := domain.New(r.Name, r.TLD)
	if err != nil {
		log.Printf("ERROR: Record.Save:domain.New: %s", err)
		return
	}

	rr := record.Record{
		Domain:     domain,
		Args:       args,
		RecordType: recordType,
		Hash:       hash,
	}

	log.Printf("%+v", rr)
}

// This function assumes the following:
// <name> <ttl?> <type> <args>
func (p Parser) getRecord(fields []string) (Record, error) {
	record := Record{}
	record.Name = fields[0]
	if !strings.HasSuffix(record.Name, ".") {
		record.Name += "." + p.origin
	}
	record.TTL = p.ttl
	typeIdx := 1
	// strip \sin\s
	if len(fields) > 3 {
		fields = filterIN(fields)
	}
	if len(fields) > 3 {
		ttl, err := strconv.ParseUint(fields[1], 10, 0)
		if err == nil {
			typeIdx = 2
			record.TTL = uint(ttl)
		} else {
			// detect and fix these, maybe go to their own channel/table
			log.Printf("WARN: getRecord:parseTTL: len(fields) == %d, fields: %s", len(fields), fields)
		}
	}
	if len(fields) <= typeIdx {
		log.Printf("ERROR: getRecord:setTypeIdx: len(fields) == %d, fields: %s", len(fields), fields)
		return record, errors.New("Unable to set typeIdx in getRecord.")
	}
	var err error
	record.RecordType, err = detectRecordType(fields[typeIdx])
	if err != nil {
		log.Printf("ERROR: getRecord.parseType: len(fields) == %d, fields: %s", len(fields), fields)
		return record, errors.New("Unable to detect Record Type.")
	}
	record.Args = fields[typeIdx+1:]
	record.TLD = p.tld
	return record, nil
}
