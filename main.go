package main

import (
	"github.com/jawr/dns/crawler"
	"github.com/jawr/dns/database/models/tlds"
	"github.com/jawr/dns/database/models/users"
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/rest"
	"github.com/jawr/dns/watcher"
	zonefile "github.com/jawr/dns/zonefile/parser"
	"github.com/stathat/jconfig"
	"net/http"
)

func main() {
	//go parseZonefiles()
	//go crawl()
	go func() {
		w, err := watcher.New()
		if err != nil {
			log.Error("Unable to start watcher: %s", err)
			return
		}
		w.Start()
	}()
	u, err := users.New("jess@lawrence.pm", "foobar")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("user: %+v", u)
	log.Info("passchk1: %s", users.CheckPassword(u.Email, "1234"))
	log.Info("passchk2: %s", users.CheckPassword(u.Email, "foobar"))
	log.Info("passchk3: %s", users.CheckPassword(u.Email+"2", "foobar"))
	log.Info("passchk4: %s", users.CheckPassword(u.Email, "1fff234"))
	u, err = users.New("jess@lawrence.pm", "barrfooo")
	if err != nil {
		log.Error("%s", err)
	}
	u, err = users.New("jess@lawrence.pm", "foobar")
	if err != nil {
		log.Error("%s", err)
	}
	startREST()
}

func crawl() {
	crawler := crawler.New(2)
	crawler.Start()
}

func testDetectDomainAndTLD() {
	s, t, err := tlds.DetectDomainAndTLD("ns1.google.co.uk")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("%+v", t)
	log.Info("%+v", s)
	s, t, err = tlds.DetectDomainAndTLD("ns1.google.co.ng")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("%+v", t)
	log.Info("%+v", s)
	s, t, err = tlds.DetectDomainAndTLD("ns1.google.co.com")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("%+v", t)
	log.Info("%+v", s)
	s, t, err = tlds.DetectDomainAndTLD("ns1.google.foo")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("%+v", t)
	log.Info("%+v", s)

}

func startREST() {
	h := rest.Setup()
	http.ListenAndServe(":8080", h)
}

func parseZonefiles() {
	config := jconfig.LoadConfig("config.json")
	dir := config.GetString("zonefile_dir")
	p := zonefile.New()
	files := []string{
		//"20141113-net.zone.gz",
		"20140621-biz.zone.gz",
		"20140622-biz.zone.gz",
		"20141210-biz.zone.gz",
	}
	for _, f := range files {
		err := p.SetupGunzipFile(dir + f)
		if err != nil {
			log.Error("Unable to setup %s: %s", f, err)
			return
		}
		p.Parse()
	}
}
