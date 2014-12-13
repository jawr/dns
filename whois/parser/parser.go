package parser

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"os/exec"
)

type Parser struct {
}

func New() Parser {
	parser := Parser{}
	return parser
}

func (p *Parser) Parse(d domain.Domain) error {
	defer util.Un(util.Trace())
	log.Info("Parse " + d.String())
	out, err := exec.Command("pwhois", "-j", d.String()).Output()
	if err != nil {
		return err
	}
	w := whois.New(d, out)
	return w.Insert()
}
