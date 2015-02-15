package watcher

import (
	db "github.com/jawr/dns/database/models/watchers"
	"github.com/jawr/dns/database/models/watchers/intervals"
	digParser "github.com/jawr/dns/dig/parser"
	"github.com/jawr/dns/log"
	whoisParser "github.com/jawr/dns/whois/parser"
	"github.com/robfig/cron"
)

type Watcher struct {
	Intervals   []intervals.Interval
	Cron        *cron.Cron
	DigParser   digParser.Parser
	WhoisParser whoisParser.Parser
}

func New() (Watcher, error) {
	w := Watcher{}
	intervalList, err := intervals.GetAll().List()
	if err != nil {
		return w, err
	}
	w.Intervals = intervalList
	w.Cron = cron.New()
	w.DigParser = digParser.New()
	w.WhoisParser = whoisParser.New()
	return w, nil
}

func (w Watcher) Start() {
	for _, i := range w.Intervals {
		j := i
		w.Cron.AddFunc(i.Value, func() {
			w.handler(j)
		})
	}
	w.Cron.Start()
}

func (w Watcher) handler(i intervals.Interval) {
	list, err := db.GetByInterval(i).List()
	if err != nil {
		log.Error("Error parsing interval (%d): %s", i.ID, err)
		return
	}
	log.Info("Watcher Handler: %s (%d domains)", i.Value, len(list))
	for _, watch := range list {
		w.DigParser.Exec(watch.Domain)
		w.WhoisParser.Exec(watch.Domain)
		watch.Save()
	}
}
