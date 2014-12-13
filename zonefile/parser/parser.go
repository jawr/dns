package parser

import (
	"bufio"
	"fmt"
	"github.com/jawr/dns/database/bulk"
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
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

func (p Parser) String() string {
	return fmt.Sprintf("%s (%s)", p.tld.Name, p.date)
}

func (p *Parser) Close() {
	log.Debug("Closing %s Parser", p.String())
	log.Debug("Closing %s Parser:setupFileDefer", p.String())
	p.setupFileDefer()
	log.Debug("Closing %s Parser:recordInsert", p.String())
	p.recordInsert.Close()
	log.Debug("Closing %s Parser:domainInsert", p.String())
	p.domainInsert.Close()
	log.Debug("Closed %s Parser", p.String())
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
	log.Info("Parsing Zonefile: %s", p.String())
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
		log.Info("Parse %s complete. Proceed with sql operations.", p.String())
	}()
	// insert our domains and commit our tx to avoid
	p.domainInsert.Insert()
	p.recordInsert.Insert()
	// TODO: drop index to record__%d_%d
	// create index's
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
	if err != nil {
		return err
	}
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
			return err
		}
	}
	p.domainInsert.Finish()
	// TODO: add index to record__%d_%d
	p.recordInsert.Finish()
	log.Info("Parse %s complete", p.String())
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
				log.Warn("handleVariable:$ttl: %s", err)
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
		log.Warn("handleLine:getRecord: %s", err)
		log.Warn("handleLine:line: %s", line)
		return
	}
	err = p.recordInsert.Add(&rr)
	if err != nil {
		log.Error("handleLine:recordInsert.Add: %s", err)
		return
	}
	p.recordTypes[rr.RecordType.Name] = rr.RecordType.ID
}
