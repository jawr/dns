package whois

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"time"
)

type JSON []byte

type Whois struct {
	ID     int32           `json:"id"`
	Domain domain.Domain   `json:"domain"`
	Data   json.RawMessage `json:"data"`
	Added  time.Time       `json:"added"`
}

func New(d domain.Domain, data []byte) Whois {
	return Whois{
		Domain: d,
		Data:   data,
	}
}

func (w *Whois) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO whois (domain, data) VALUES ($1, $2)",
		w.Domain.UUID.String(),
		string(w.Data),
	)
	return err
}
