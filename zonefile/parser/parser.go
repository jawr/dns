package parser

import (
	"bufio"
	"fmt"
	"github.com/jawr/dns/database/bulk"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records"
	db "github.com/jawr/dns/database/models/zonefile/parser"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"strconv"
	"strings"
)

type Parser struct {
	scanner        *bufio.Scanner
	setupFileDefer func()
	ttl            uint
	origin         string
	originCheck    bool
	line           string
	lineCount      uint
	domainInsert   *bulk.Insert
	recordInsert   *bulk.Insert
	recordTypes    map[string]int32
	db.Parser
}

func New() Parser {
	parser := Parser{
		ttl:         86400, //24 hours
		originCheck: false,
	}
	return parser
}

func (p *Parser) SetOrigin(s string) {
	p.origin = s
}

func (p Parser) String() string {
	return fmt.Sprintf("%s (%s)", p.TLD.Name, p.Date.Format("02-01-2006"))
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
	p.Update("Parse finished.")
	p.Finish()
}

func (p *Parser) Parse() error {
	defer p.Close()
	defer util.Un(util.Trace())
	ri, err := records.NewBulkInsert()
	if err != nil {
		return err
	}
	p.recordInsert = &ri
	bi, err := domains.NewBulkInsert()
	if err != nil {
		return err
	}
	p.recordTypes = make(map[string]int32)
	p.domainInsert = &bi
	log.Info("Parsing Zonefile: %s", p.String())
	err = p.Update("Load file in to temporary tables.")
	if err != nil {
		log.Error("%s", err)
		return err
	}
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
		err = p.Update("File read in to temporary tables.")
		if err != nil {
			log.Error("Unable to read in temporary tables: %s", err)
			return
		}
		log.Info("Parse %s complete. Proceed with sql operations.", p.String())
	}()
	// insert our domains and commit our tx to avoid
	err = p.Update("Close insert domains statement.")
	if err != nil {
		return err
	}
	p.domainInsert.Insert()
	err = p.Update("Close insert records statement.")
	if err != nil {
		return err
	}
	p.recordInsert.Insert()
	// TODO: drop index to record__%d_%d
	// create index's
	err = p.Update("Create temp domain table index.")
	if err != nil {
		return err
	}

	err = p.domainInsert.Index(fmt.Sprintf("CREATE INDEX _domain__%d_uuid_idx ON %%s (uuid)", p.TLD.ID))
	if err != nil {
		return err
	}
	err = p.Update("Create temp record table index.")
	err = p.recordInsert.Index("CREATE INDEX uuid_idx ON %s (uuid)")
	if err != nil {
		return err
	}
	err = p.Update("Merge temp domain with domain.")
	err = p.domainInsert.Merge(
		fmt.Sprintf(`
			INSERT INTO domain__%d
			SELECT DISTINCT * FROM %%s d2
				WHERE NOT EXISTS (
					SELECT NULL FROM domain__%d d WHERE
						d.uuid = d2.uuid
				)`,
			p.TLD.ID,
			p.TLD.ID,
		),
	)
	if err != nil {
		return err
	}
	for rtName, rtID := range p.recordTypes {
		err = p.recordInsert.Exec(fmt.Sprintf(
			"DROP INDEX IF EXISTS record__%d_%d_idx",
			rtID,
			p.TLD.ID,
		))
		if err != nil {
			return err
		}
		err = p.Update("Merge temp record with record (type: %s).", rtName)
		if err != nil {
			return err
		}
		err = p.recordInsert.Merge(
			fmt.Sprintf(`
				INSERT INTO record__%d_%d
				SELECT DISTINCT ON (uuid) *, %d FROM %%s r2
					WHERE NOT EXISTS (
						SELECT NULL FROM record__%d_%d r WHERE
							r.uuid = r2.uuid AND r.record_type = %d
					) AND r2.record_type = %d`,
				rtID,
				p.TLD.ID,
				p.ID,
				rtID,
				p.TLD.ID,
				rtID,
				rtID,
			),
		)
		if err != nil {
			return err
		}
		err = p.recordInsert.Exec(fmt.Sprintf(
			"CREATE INDEX record__%d_%d_idx ON record__%d_%d USING HASH (domain)",
			rtID,
			p.TLD.ID,
			rtID,
			p.TLD.ID,
		))
		if err != nil {
			return err
		}
	}
	err = p.Update("Commit domain insert.")
	if err != nil {
		return err
	}
	p.domainInsert.Finish()
	// TODO: add index to record__%d_%d
	err = p.Update("Commit record insert.")
	if err != nil {
		return err
	}
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
	rr, err := records.New(line, p.origin, p.TLD, p.ttl, p.Date)
	if err != nil {
		log.Warn("handleLine:getRecord: %s", err)
		log.Warn("handleLine:line: %s", line)
		return
	}
	err = p.domainInsert.Add(&rr.Domain)
	if err != nil {
		log.Error("handleLine: Unable to bulk insert Domain: %s", err)
		return
	}
	err = p.recordInsert.Add(&rr)
	if err != nil {
		log.Error("handleLine:recordInsert.Add: %s", err)
		return
	}
	p.recordTypes[rr.Type.Name] = rr.Type.ID
}
