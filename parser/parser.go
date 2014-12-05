package parser

import (
	"bufio"
	"compress/gzip"
	"errors"
	"github.com/jawr/dns/database/tld"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	scanner *bufio.Scanner
	tld     *tld.TLD
	ttl     uint
	origin  string
}

func New(tldName string) Parser {
	tld := tld.TLD{
		ID:   1,
		Name: tldName,
	}
	parser := Parser{
		tld:    &tld,
		ttl:    86400, //24 hours
		origin: tld.Name + ".",
	}
	return parser
}

func (p *Parser) SetupGunzipFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	p.scanner = bufio.NewScanner(reader)
	return nil
}

func (p *Parser) Parse() error {
	log.Println("Starting parse")
	count := 0
	var previous string
	for p.scanner.Scan() {
		if count > 200 {
			//break
		}
		count++
		line := strings.ToLower(p.scanner.Text())
		if len(line) == 0 {
			continue
		}
		commentIdx := strings.Index(line, ";")
		if commentIdx > 0 {
			//comment := line[commentIdx:]
			line = line[:commentIdx]
		}
		switch line[0] {
		case ';':
			p.handleComment(line)
			break
		case '$':
			p.handleVariable(line)
			break
		case ' ':
			p.handleZonedLine(line, previous)
			break
		default:
			p.handleLine(line)
		}
		previous = line
	}
	return nil
}

func (p *Parser) handleComment(line string) {
}

func (p *Parser) handleVariable(line string) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		switch fields[0] {
		case "$origin":
			p.origin = fields[1]
			break
		case "$ttl":
			ttl, err := strconv.ParseUint(fields[1], 10, 0)
			if err != nil {
				log.Printf("Error parsing Variable 'ttl': %s", err)
				return
			}
			p.ttl = uint(ttl)
			break
		}
	}
}

func (p *Parser) handleZonedLine(line, previous string) {
	log.Printf("Zoned line: %s | %s", previous, line)
}

var classTypeRe *regexp.Regexp = regexp.MustCompile(`[\s\t](cname|mx|a|aaaa|ns)[\s\t]`)

//var classTypeRe *regexp.Regexp = regexp.MustCompile(`[\s\t](a)[\s\t]`)

func (p *Parser) handleLine(line string) {
	if !classTypeRe.MatchString(line) {
		return
	}
	fields := strings.Fields(line)
	log.Println(fields)
	record, err := p.getRecord(fields)
	if err != nil {
		return
	}
	log.Println(record)
}

func parseARecord(name, recordType, addr string, ttl uint) {

}

type Record struct {
	Name       string
	Args       []string
	TTL        uint
	RecordType string
}

var allClassTypesRe *regexp.Regexp = regexp.MustCompile(`(afsdb|apl|caa|cert|dhcid|dlv|dname|dnskey|ds|hip|ipseckey|key|kx|loc|naptr|nsec|nsec3|nsec3param|ptr|rrsig|rp|sig|soa|spf|srv|sshfp|ta|tkey|tlsa|tsig|txt|axfr|ixfr|opt|cname|mx|a|aaaa|ns)`)

/*
	this function assumes the following:

	<name> <ttl?> <type> <args>
*/
func (p Parser) getRecord(fields []string) (Record, error) {
	record := Record{}
	record.Name = fields[0]
	if !strings.HasSuffix(record.Name, ".") {
		record.Name += "." + p.origin
	}
	record.TTL = p.ttl
	typeIdx := 1
	if len(fields) > 3 {
		ttl, err := strconv.ParseUint(fields[1], 10, 0)
		if err == nil {
			typeIdx = 2
			record.TTL = uint(ttl)
		} else {
			// detect and fix these, maybe go to their own channel/table
			log.Printf("ERROR: getRecord:parseTTL: len(fields) == %d, fields: %s", len(fields), fields)
		}
	}
	if !allClassTypesRe.MatchString(fields[typeIdx]) {
		log.Printf("ERROR: getRecord.parseType: len(fields) == %d, fields: %s", len(fields), fields)
		return record, errors.New("Unable to detect Record Type.")
	}
	record.RecordType = fields[typeIdx]
	record.Args = fields[typeIdx+1:]
	return record, nil
}
