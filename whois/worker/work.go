package worker

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/log"
)

type Request struct {
	Query string
}

func (r Request) Do(w Worker) {
	s, t, err := tld.DetectDomainAndTLD(r.Query)
	if err != nil {
		log.Error("worker%d: unable to detect TLD and domain: %s", w.ID, err)
		return
	}
	d, err := domain.Get(domain.GetByNameAndTLD(), s, t.ID)
	if err != nil {
		log.Error("worker%d: unable to get domain: %s", w.ID, err)
		return
	}
	_, err = w.Parser.Exec(d)
	if err != nil {
		log.Error("worker%d: unable to parse domain: %s", w.ID, err)
		return
	}
	log.Info("worker%d: got result for domain", w.ID)
}
