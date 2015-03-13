package watcher

import (
	"github.com/jawr/dns/database/models/notifications"
	db "github.com/jawr/dns/database/models/watchers"
	"github.com/jawr/dns/database/models/watchers/intervals"
	"github.com/jawr/dns/database/models/whois"
	digParser "github.com/jawr/dns/dig/parser"
	"github.com/jawr/dns/log"
	whoisParser "github.com/jawr/dns/whois/parser"
	"github.com/robfig/cron"
	"time"
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
	for _, watch := range list {
		w.DigParser.Exec(watch.Domain)
		before, err := whois.GetByDomain(watch.Domain).Count()
		if err != nil {
			log.Error("%s", err)
		}
		w.WhoisParser.Exec(watch.Domain)
		after, err := whois.GetByDomain(watch.Domain).Count()
		if err != nil {
			log.Error("%s", err)
		}
		if after > before && before > 0 {
			log.Info("%s whois changed", watch.Domain)
			message := notifications.Message{
				Added:   time.Now(),
				Message: "Domain Whois update.",
				Domain:  watch.Domain,
			}
			for _, user := range watch.Users {
				log.Info("%s", user.Email)
				noter, _ := notifications.SetupUserNotification(user)
				noter.AddMessage(message)
			}
		}
		err = watch.Save()
		if err != nil {
			log.Error("%s", err)
		}

	}
}
