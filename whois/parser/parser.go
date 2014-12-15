package parser

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/util"
	"github.com/sbinet/go-python"
	"os/exec"
	"sync"
)

type Parser struct {
}

func New() Parser {
	parser := Parser{}
	return parser
}

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
}

var lock sync.Mutex

func (p *Parser) Parse(d domain.Domain) (whois.Result, error) {
	defer util.Un(util.Trace())
	lock.Lock()
	defer lock.Unlock()
	pwhois := python.PyImport_ImportModuleNoBlock("pythonwhois.net")
	getResult := pwhois.GetAttrString("get_whois_raw")

	//log.Info("%+v", getResult)

	args := python.PyTuple_New(1)
	python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString(d.String()))
	result := getResult.CallObject(args)
	log.Info("%s", result)
	log.Info("%s", python.PyString_AsString(result))
	log.Info("%s", python.PyByteArray_AsString(result))

	//w := whois.New(d, out)
	//return w, w.Insert()
	return whois.Result{}, nil
}

func (p *Parser) Exec(d domain.Domain) (whois.Result, error) {
	defer util.Un(util.Trace())
	log.Info("Parse " + d.String())
	out, err := exec.Command("pwhois", "-j", d.String()).Output()
	if err != nil {
		return whois.Result{}, err
	}
	w := whois.New(d, out)
	return w, w.Insert()
}
