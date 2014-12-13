package parser

import (
	"bufio"
	"fmt"
	"github.com/jawr/dns/database/bulk"
	"github.com/jawr/dns/database/domain"
	"github.com/jawr/dns/database/record"
	"github.com/jawr/dns/database/tld"
	"github.com/jawr/dns/util"
	"log"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	scanner        *bufio.Scanner
	setupFileDefer func()
	tld            tld.TLD
	tldName        string
	date           time.Time
	ttl            uint
	origin         string
	originCheck    bool
	line           string
	lineCount      uint
	domainInsert   *bulk.Insert
	recordInsert   *bulk.Insert
	recordTypes    map[string]int32
}

func New() Parser {
	parser := Parser{
		ttl:         86400, //24 hours
		originCheck: false,
	}
	return parser
}

func (p *Parser) Close() {
	log.Printf("INFO: Closing %s Parser", p.tld.Name)
	log.Printf("INFO: Closing %s Parser:setupFileDefer", p.tld.Name)
	p.setupFileDefer()
	log.Printf("INFO: Closing %s Parser:recordInsert", p.tld.Name)
	p.recordInsert.Close()
	log.Printf("INFO: Closing %s Parser:domainInsert", p.tld.Name)
	p.domainInsert.Close()
	log.Printf("INFO: Closed %s Parser", p.tld.Name)
}

func (p *Parser) Parse() error {
	defer p.Close()
	defer util.Un(util.Trace())
	ri, err := record.NewBulkInsert()
	if err != nil {
		return err
	}
	p.recordInsert = &ri
	bi, err := domain.NewBulkInsert()
	if err != nil {
		return err
	}
	p.recordTypes = make(map[string]int32)
	p.domainInsert = &bi
	log.Printf("INFO: Parsing %s zonefile", p.tld.Name)
	var previous string
	p.lineCount = 0
	func() {
		defer util.Un(util.Trace())
		for p.scanner.Scan() {
			p.lineCount++
			p.line = strings.ToLower(p.scanner.Text())
			line := p.line
			commentIdx := strings.Index(line, ";")
			if commentIdx > 0 {
				//comment := line[commentIdx:]
				line = line[:commentIdx]
			}
			if len(line) == 0 {
				continue
			}
			switch line[0] {
			case ';':
				p.handleComment(line)
			case '@':
				// need to wdd more ways of detecting SOA, could have a switch
				// that only goes to handleLine when it has parsed $origin var
				p.handleSOA(line)
			case '$':
				p.handleVariable(line)
			case ' ':
			case '	':
				p.handleZonedLine(line, previous)
			default:
				if !p.originCheck {
					p.handleSOA(line)
				} else {
					p.handleLine(line)
				}
			}
			previous = line
		}
		log.Println("INFO: Parse file complete. Proceed with sql operations.")
	}()
	// insert our domains and commit our tx to avoid
	p.domainInsert.Insert()
	p.recordInsert.Insert()
	err = p.domainInsert.Index("CREATE INDEX uuid_idx ON %s (uuid)")
	err = p.recordInsert.Index("CREATE INDEX uuid_idx ON %s (uuid)")
	err = p.domainInsert.Merge(
		fmt.Sprintf(`
			INSERT INTO domain__%d
			SELECT DISTINCT * FROM %%s d2
				WHERE NOT EXISTS (
					SELECT NULL FROM domain__%d d WHERE
						d.uuid = d2.uuid
				)`,
			p.tld.ID,
			p.tld.ID,
		),
	)
	//err = p.domainInsert.Merge("INSERT INTO domain SELECT * FROM %s")
	if err != nil {
		log.Println("ERROR: Parse:domain Merge: %s", err)
		return err
	}
	//err = p.recordInsert.Merge("INSERT INTO record SELECT * FROM %s")
	// this one might still be slow due to us not inserting straight in to the
	// correct record type table (bypass the trigger)
	for _, rtID := range p.recordTypes {
		err = p.recordInsert.Merge(
			fmt.Sprintf(`
				INSERT INTO record__%d_%d
				SELECT DISTINCT ON (uuid) * FROM %%s r2
					WHERE NOT EXISTS (
						SELECT NULL FROM record__%d_%d r WHERE
							r.uuid = r2.uuid AND r.record_type = %d
					) AND r2.record_type = %d`,
				rtID,
				p.tld.ID,
				rtID,
				p.tld.ID,
				rtID,
				rtID,
			),
		)
		if err != nil {
			log.Println("ERROR: Parse:Merge: %s", err)
			return err
		}
	}
	p.domainInsert.Finish()
	p.recordInsert.Finish()
	log.Println("INFO: Parse complete")
	return nil
}

func (p *Parser) handleSOA(line string) {
}

func (p *Parser) handleComment(line string) {
}

func (p *Parser) handleZonedLine(line, previous string) {
}

func (p *Parser) handleVariable(line string) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		switch fields[0] {
		case "$origin":
			p.origin = fields[1]
			p.originCheck = true
		case "$ttl":
			ttl, err := strconv.ParseUint(fields[1], 10, 0)
			if err != nil {
				log.Printf("WARN: handleVariable:$ttl: %s", err)
				return
			}
			p.ttl = uint(ttl)
		}
	}
}

func (p *Parser) handleLine(line string) {
	fields := strings.Fields(line)
	rr, err := p.getRecord(fields)
	if err != nil {
		log.Printf("WARN: handleLine:getRecord: %s", err)
		log.Printf("WARN: handleLine:line: %s", line)
		return
	}
	err = p.recordInsert.Add(&rr)
	if err != nil {
		log.Printf("ERROR: handleLine:recordInsert.Add: %s", err)
		return
	}
	p.recordTypes[rr.RecordType.Name] = rr.RecordType.ID
}
