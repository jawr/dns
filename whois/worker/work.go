package worker

import (
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/whois"
	"github.com/jawr/dns/log"
)

type Request struct {
	Domain domains.Domain
	Record chan whois.Record
}

func (r Request) Do(w Worker) {
	res, err := w.Parser.Exec(r.Domain)
	if err != nil {
		log.Error("worker%d: unable to parse domain: %s", w.ID, err)
		return
	}
	select {
	case r.Record <- res:
		break
	default:
		break
	}
}
