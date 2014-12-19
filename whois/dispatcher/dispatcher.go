package dispatcher

import (
	"github.com/jawr/dns/whois/worker"
)

var WorkerQueue chan chan worker.Request

type Dispatcher struct {
	Workers chan chan worker.Request
	Work    chan worker.Request
}

func New(nworkers int) Dispatcher {
	d := Dispatcher{
		Workers: make(chan chan worker.Request, nworkers),
		Work:    make(chan worker.Request, nworkers*10),
	}

	for i := 0; i < nworkers; i++ {
		w := worker.New(i+1, d.Workers)
		w.Start()
	}

	return d
}

func (d Dispatcher) Start() {
	go func() {
		for {
			select {
			case work := <-d.Work:
				go func() {
					w := <-d.Workers
					w <- work
				}()
			}
		}
	}()
}
