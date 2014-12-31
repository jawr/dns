package parser

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
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

func (p *Parser) Exec(d domain.Domain) (whois.Result, error) {
	defer util.Un(util.Trace())
	log.Info("Parse " + d.String())
	out, err := exec.Command("pwhois", "-j", d.String()).Output()
	if err != nil {
		return whois.Result{}, err
	}
	w, err := whois.New(d, out)
	if err != nil {
		log.Error("Parse Whois for %s: %s", d.String(), err)
		return w, err
	}
	return w, w.Insert()
}
