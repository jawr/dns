package parser

import (
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	"os/exec"
	"sync"
)

type Parser struct {
}

func New() Parser {
	parser := Parser{}
	return parser
}

var lock sync.Mutex

func (p *Parser) Exec(d domains.Domain) (whois.Record, error) {
	out, err := exec.Command("pwhois", "-j", d.String()).Output()
	if err != nil {
		log.Error("Whois error for " + d.String())
		return whois.Record{}, err
	}
	w, err := whois.New(d, out)
	if err != nil {
		log.Error("Parse Whois for %s: %s", d.String(), err)
		return w, err
	}
	return w, w.Insert()
}
