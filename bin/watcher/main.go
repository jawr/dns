package main

import (
	"github.com/jawr/dns/log"
	"github.com/jawr/dns/watcher"
)

func main() {
	w, err := watcher.New()
	if err != nil {
		log.Error("Unable to start Watcher: %s", err)
		return
	}
	w.Start()
	q := make(chan bool)
	<-q
}
