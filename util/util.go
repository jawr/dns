package util

import (
	"github.com/jawr/dns/log"
	"time"
)

func Trace() time.Time {
	return time.Now()
}

func Un(t time.Time) {
	elapsed := time.Since(t)
	log.Debug("%s", elapsed.String())
}

func FilterIN(sl []string) []string {
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
