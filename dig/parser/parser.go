package parser

import (
	//"os/exec"
	"github.com/jawr/dns/database/models/dig"
)

type Parser struct {
}

func New() Parser {
	parser := Parser{}
	return parser
}

func (p *Parser) Parse(d domain.Domain) (dig.Result, error) {
	d := dig.Result{}
	return d, nil
}
