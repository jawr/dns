package crawler

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/dig/dispatcher"
	"github.com/jawr/dns/log"
	"time"
)

type Crawler struct {
	quit   chan bool
	delay  time.Duration
	offset int
}

const WINDOW int = 1000

func New(delay int) Crawler {
	return Crawler{
		delay:  time.Duration(delay) * time.Second,
		quit:   make(chan bool),
		offset: 0,
	}
}

func (c Crawler) Start() {
	go func() {
		for {
			select {
			case <-c.quit:
				log.Info("Quit whois crawler.")
				return
			default:
				domains, err := domain.GetList(domain.GetAllLimitOffset(), WINDOW, c.offset)
				if err != nil {
					log.Error("Whois crawler. Unable to get domains: %s", err)
					// shutdown
					c.Stop()
					return
				}
				for _, d := range domains {
					dispatcher.AddDomain(d)
					time.Sleep(c.delay)
				}
			}
			c.offset += WINDOW
		}
	}()
}

func (c Crawler) Stop() {
	go func() {
		c.quit <- true
	}()
}
