package record

import (
	"github.com/jawr/dns/database/domain"
	"github.com/jawr/dns/database/record_type"
	"time"
)

type RecordArgs struct {
	TTL  uint
	Addr string
	Args []string
}

type Record struct {
	ID         uint
	Domain     domain.Domain
	Name       string
	Hash       string
	Args       RecordArgs
	RecordType record_type.RecordType
	Added      *time.Time
	Active     bool
	Hitory     byte
}
