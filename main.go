package main

import (
	"github.com/jawr/dns/database/tld"
	"github.com/jawr/dns/parser"
	"io"
	"log"
	"os"
)

func main() {
	// ofload to another function, but how to deal with defer close??
	f, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR: SetupLog:OpenFile: %s", err)
		return
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	t, err := tld.New("biz")
	if err != nil {
		log.Println(err)
		return
	}

	p, err := parser.New(t)
	if err != nil {
		log.Printf("Error setting up Parser: %s", err)
		return
	}
	err = p.SetupGunzipFile("/home/jawr/dns/zonefiles/20141113-net.zone.gz")
	//err = p.SetupFile("/home/jawr/dns/zonefiles/biz.zone")
	//err = p.SetupFile("/home/jawr/dns/zonefiles/20140622-biz.zone")
	if err != nil {
		log.Printf("Error opening Gunzip file for parsing: %s", err)
		return
	}
	p.Parse()
}

// 2014/12/05 00:47:07 Starting parse
// 2014/12/05 01:03:31 {kitpvp.deadlystars.net. [104.218.96.198] 172800 a}
