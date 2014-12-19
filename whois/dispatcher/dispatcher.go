package dispatcher

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/database/models/tld"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/whois/worker"
)

var Workers chan chan worker.Request
var Work chan worker.Request

func init() {
	Start(2)
}

func Start(nworkers int) {
	Workers = make(chan chan worker.Request, nworkers)
	Work = make(chan worker.Request, nworkers*10)

	for i := 0; i < nworkers; i++ {
		w := worker.New(i+1, Workers)
		w.Start()
	}

	go func() {
		for {
			select {
			case work := <-Work:
				go func() {
					w := <-Workers
					w <- work
				}()
			}
		}
	}()
}

func AddDomain(d domain.Domain) {
	wr := worker.Request{Domain: d}
	Work <- wr
}

func AddQuery(q string) {
	// could offload this in to a seperate anon function to avoid bottleneck
	s, t, err := tld.DetectDomainAndTLD(q)
	if err != nil {
		log.Error("Whois dipatcher: Unable to detect TLD and domain: %s (%s)", err, q)
		return
	}
	d, err := domain.Get(domain.GetByNameAndTLD(), s, t.ID)
	if err != nil {
		log.Error("Whois dispatcher: unable to get domain: %s (%s)", err, s)
		return
	}
	wr := worker.Request{Domain: d}
	Work <- wr
}

// TODO: add a stop/quit
