package parser

import (
	"log"
	"time"
)

func trace() time.Time {
	return time.Now()
}

func un(t time.Time) {
	elapsed := time.Since(t)
	log.Printf("%s", elapsed.String())
}
