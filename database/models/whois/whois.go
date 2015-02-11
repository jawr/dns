package whois

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"regexp"
	"time"
)

var uuidRawRE *regexp.Regexp = regexp.MustCompile(`([0-1]?\d|2[0-3]):([0-5]?\d):([0-5]?\d)`)

type JSON []byte

type Record struct {
	ID       int32           `json:"id"`
	Domain   domains.Domain  `json:"domain"`
	Data     json.RawMessage `json:"data"`
	Raw      json.RawMessage `json:"raw"`
	Contacts json.RawMessage `json:"contacts"`
	Emails   []string        `json:"emails"`
	Added    time.Time       `json:"added"`
	UUID     uuid.UUID       `json:"uuid"`
}

func New(d domains.Domain, data []byte) (Record, error) {
	raw := Raw{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		panic(err)
	}

	emails := parseEmails(raw)

	rawContacts, err := json.Marshal(&raw.Contacts)
	if err != nil {
		panic(err)
	}
	rawData, err := json.Marshal(&raw.Data)
	uuidRaw := fmt.Sprintf("%s", raw.Raw)
	uuidRaw = uuidRawRE.ReplaceAllString(uuidRaw, "")
	id := uuid.NewSHA1(uuid.NameSpace_OID, []byte(uuidRaw))
	return Record{
		Domain:   d,
		Data:     rawData,
		Raw:      raw.Raw,
		Contacts: rawContacts,
		Emails:   emails,
		UUID:     id,
	}, err
}

func (r Record) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	emails, err := json.Marshal(r.Emails)
	if err != nil {
		return err
	}
	data, err := r.Data.MarshalJSON()
	if err != nil {
		return err
	}
	raw, err := r.Raw.MarshalJSON()
	if err != nil {
		return err
	}
	contacts, err := r.Contacts.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = conn.Exec(`INSERT INTO whois 
					(domain, data, raw_whois, contacts, emails, uuid) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
		r.Domain.UUID.String(),
		string(data),
		string(raw),
		string(contacts),
		string(emails),
		r.UUID.String(),
	)
	return err
}
