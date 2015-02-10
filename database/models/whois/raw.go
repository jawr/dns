package whois

import (
	"encoding/json"
)

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
	City         string `json:"city,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	State        string `json:"state,omitempty"`
	Handle       string `json:"handle,omitempty"`
	Street       string `json:"street,omitempty"`
	Country      string `json:"country,omitempty"`
	Facsimile    string `json:"facsimilie,omitempty"`
	Postalcode   string `json:"postalcode,omitempty"`
	Organization string `json:"organization,omitempty"`
}

type Contacts struct {
	Tech       *Contact `json:"tech,omitempty"`
	Admin      *Contact `json:"admin,omitempty"`
	Billing    *Contact `json:"billing,omitempty"`
	Registrant *Contact `json:"registrant,omitempty"`
}

type Raw struct {
	Data
	Raw      json.RawMessage `json:"raw"`
	Emails   []string        `json:"emails"`
	Contacts Contacts        `json:"contacts"`
}

func parseEmails(raw Raw) []string {
	if len(raw.Emails) == 0 {
		return []string{}
	}
	emails := make(map[string]bool)
	for _, v := range raw.Emails {
		emails[v] = true
	}
	if raw.Contacts.Tech != nil {
		emails[raw.Contacts.Tech.Email] = true
	}
	if raw.Contacts.Admin != nil {
		emails[raw.Contacts.Admin.Email] = true
	}
	if raw.Contacts.Billing != nil {
		emails[raw.Contacts.Billing.Email] = true
	}
	if raw.Contacts.Registrant != nil {
		emails[raw.Contacts.Registrant.Email] = true
	}

	var list = make([]string, 0)
	for k, _ := range emails {
		if k != "" {
			list = append(list, k)
		}
	}

	return list
}
