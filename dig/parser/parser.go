package parser

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"os/exec"
	"strings"
	"time"
)

type Parser struct {
}

func New() Parser {
	return Parser{}
}

func (p *Parser) Exec(d domain.Domain) ([]record.Record, error) {
	defer util.Un(util.Trace())
	log.Info("Parse Dig " + d.String())
	out, err := exec.Command("dig", "all", d.String()).Output()
	if err != nil {
		return []record.Record{}, err
	}
	date := time.Now()
	origin := d.TLD.Name + "."
	lines := strings.Split(strings.ToLower(string(out)), "\n")
	results := make([]record.Record, 0)
	for _, line := range lines {
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		rr, err := record.New(line, origin, d.TLD, 86400, date)
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
