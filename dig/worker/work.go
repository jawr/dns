package worker

import (
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/records"
	"github.com/jawr/dns/log"
)

type Request struct {
	Domain domains.Domain
	Result chan []records.Record
}

func (r Request) Do(w Worker) {
	res, err := w.Parser.Exec(r.Domain)
	if err != nil {
		log.Error("worker%d: unable to parse domain: %s", w.ID, err)
		return
	}
	log.Info("res: %v", res)
	select {
	case r.Result <- res:
		break
	default:
		break
	}
}
