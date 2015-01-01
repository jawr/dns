package whois

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domain"
	"regexp"
	"time"
)

var uuidRawRE *regexp.Regexp = regexp.MustCompile(`([0-1]?\d|2[0-3]):([0-5]?\d):([0-5]?\d)`)

type JSON []byte

type Result struct {
	ID       int32           `json:"id"`
	Domain   domain.Domain   `json:"domain"`
	Data     json.RawMessage `json:"data"`
	Raw      json.RawMessage `json:"raw"`
	Contacts json.RawMessage `json:"contacts"`
	Emails   []string        `json:"emails"`
	Added    time.Time       `json:"added"`
	UUID     uuid.UUID       `json:"uuid"`
}

type Data struct {
	ID             []string `json:"id,omitempty"`
	Status         []string `json:"status,omitempty"`
	Registrar      []string `json:"registrar,omitempty"`
	Nameservers    []string `json:"nameservers,omitempty"`
	UpdatedDate    []string `json:"updated_date,omitempty"`
	WhoisServer    []string `json:"whois_server,omitempty"`
	CreationDate   []string `json:"creation_date,omitempty"`
	ExpirationDate []string `json:"expiration_date,omitempty"`
}

type Contact struct {
	City         string `json:"city"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	State        string `json:"state"`
	Handle       string `json:"handle"`
	Street       string `json:"street"`
	Country      string `json:"country"`
	Facsimile    string `json:"facsimilie"`
	Postalcode   string `json:"postalcode"`
	Organization string `json:"organization"`
}

type Contacts struct {
	Tech       Contact `json:"tech,omitempty"`
	Admin      Contact `json:"admin,omitempty"`
	Billing    Contact `json:"billing,omitempty"`
	Registrant Contact `json:"registrant,omitempty"`
}

type Raw struct {
	Data
	Raw      json.RawMessage `json:"raw"`
	Emails   []string        `json:"emails"`
	Contacts Contacts        `json:"contacts"`
}

func New(d domain.Domain, data []byte) (Result, error) {
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
	return Result{
		Domain:   d,
		Data:     rawData,
		Raw:      raw.Raw,
		Contacts: rawContacts,
		Emails:   emails,
		UUID:     id,
	}, err
}

func (w *Result) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	emails, err := json.Marshal(w.Emails)
	if err != nil {
		return err
	}
	data, err := w.Data.MarshalJSON()
	if err != nil {
		return err
	}
	raw, err := w.Raw.MarshalJSON()
	if err != nil {
		return err
	}
	contacts, err := w.Contacts.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = conn.Exec(`INSERT INTO whois 
					(domain, data, raw_whois, contacts, emails, uuid) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
		w.Domain.UUID.String(),
		string(data),
		string(raw),
		string(contacts),
		string(emails),
		w.UUID.String(),
	)
	return err
}

func parseEmails(raw Raw) []string {
	emails := make(map[string]bool)
	for _, v := range raw.Emails {
		emails[v] = true
	}
	emails[raw.Contacts.Tech.Email] = true
	emails[raw.Contacts.Admin.Email] = true
	emails[raw.Contacts.Billing.Email] = true
	emails[raw.Contacts.Registrant.Email] = true

	var list = make([]string, 0)
	for k, _ := range emails {
		if k != "" {
			list = append(list, k)
		}
	}

	return list
}
