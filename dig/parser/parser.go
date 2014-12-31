package parser

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	zonefile "github.com/jawr/dns/zonefile/parser"
	"os/exec"
	"strings"
	"time"
)

type Parser struct {
	zp zonefile.Parser
}

func New() Parser {
	zp := zonefile.New()
	zp.Date = time.Now()
	parser := Parser{
		zp: zp,
	}
	return parser
}

func (p *Parser) Exec(d domain.Domain) ([]record.Record, error) {
	defer util.Un(util.Trace())
	log.Info("Parse " + d.String())
	out, err := exec.Command("dig", "all", d.String()).Output()
	if err != nil {
		return []record.Record{}, err
	}
	p.zp.SetOrigin(d.TLD.Name + ".")
	p.zp.TLD = d.TLD
	lines := strings.Split(strings.ToLower(string(out)), "\n")
	results := make([]record.Record, 0)
	for _, line := range lines {
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		fields := strings.Fields(line)
		rr, err := p.zp.GetRecord(fields)
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
