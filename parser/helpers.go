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

func filterIN(sl []string) []string {
	fn := func(s string) bool {
		if s == "in" {
			return true
		}
		return false
	}
	return filter(sl, fn)
}

func filter(sl []string, fn func(string) bool) []string {
	outi := 0
	res := sl
	for _, v := range sl {
		if !fn(v) {
			res[outi] = v
			outi++
		}
	}
	return res[0:outi]
}
