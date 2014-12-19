package worker

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/log"
)

type Request struct {
	Domain domain.Domain
}

func (r Request) Do(w Worker) {
	_, err := w.Parser.Exec(r.Domain)
	if err != nil {
		log.Error("worker%d: unable to parse domain: %s", w.ID, err)
		return
	}
	log.Info("worker%d: got result for domain", w.ID)
}
