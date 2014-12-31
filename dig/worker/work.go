package worker

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/record"
	"github.com/jawr/dns/log"
)

type Request struct {
	Domain domain.Domain
	Result chan []record.Record
}

func (r Request) Do(w Worker) {
	res, err := w.Parser.Exec(r.Domain)
	if err != nil {
		log.Error("worker%d: unable to parse domain: %s", w.ID, err)
		return
	}
	select {
	case r.Result <- res:
		break
	default:
		break
	}
}
