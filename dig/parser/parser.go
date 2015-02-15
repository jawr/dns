package parser

import (
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records"
	"github.com/jawr/dns/log"
	"os/exec"
	"strings"
	"time"
)

type Parser struct {
}

func New() Parser {
	return Parser{}
}

func (p *Parser) Exec(d domains.Domain) ([]records.Record, error) {
	out, err := exec.Command("dig", "all", d.String()).Output()
	if err != nil {
		return []records.Record{}, err
	}
	date := time.Now()
	origin := d.TLD.Name + "."
	lines := strings.Split(strings.ToLower(string(out)), "\n")
	results := make([]records.Record, 0)
	for _, line := range lines {
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		rr, err := records.New(line, origin, d.TLD, 86400, date)
		if err != nil {
			log.Error("%s: %s", line, err)
			continue
		}
		err = rr.Insert()
		if err != nil {
			log.Error("%s: %s", line, err)
			continue
		}
	}
	return results, nil
}
