package main

import (
	"github.com/jawr/dns/database/models/domain"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest"
	whois "github.com/jawr/dns/whois/parser"
	zonefile "github.com/jawr/dns/zonefile/parser"
	"net/http"
)

func main() {
	h := rest.Setup()
	http.ListenAndServe(":8080", h)
}

func parseWhois() {
	d, err := domain.Get(domain.GetByUUID(), "df16d75f-7d86-51dd-9951-4b19e723a6d2")
	if err != nil {
		log.Error("Unable to get domain: %s", err)
		return
	}

	p := whois.New()
	err = p.Parse(d)
	if err != nil {
		log.Error("Unable to parse domain: %s", err)
	}
}

func parseZonefiles() {
	p := zonefile.New()
	files := []string{
		//"20141113-net.zone.gz",
		"20140621-biz.zone.gz",
		"20140622-biz.zone.gz",
		"20141210-biz.zone.gz",
	}
	for _, f := range files {
		err := p.SetupGunzipFile("/home/jawr/dns/zonefiles/" + f)
		if err != nil {
			log.Error("Unable to setup %s: %s", f, err)
			return
		}
		p.Parse()
	}
}
