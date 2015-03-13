package whois

import (
	"encoding/json"
	"reflect"
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

func (r *Record) parseRaw(raw Raw) {
	emails := make(map[string]bool)
	organizations := make(map[string]bool)
	phones := make(map[string]bool)
	postcodes := make(map[string]bool)
	names := make(map[string]bool)

	v := reflect.ValueOf(raw.Contacts)
	for _, v := range raw.Emails {
		emails[v] = true
	}

	for i := 0; i < v.NumField(); i++ {
		contact := v.Field(i).Interface().(*Contact)
		if contact != nil {
			cv := reflect.ValueOf(*contact)
			for j := 0; j < cv.NumField(); j++ {
				field := cv.Type().Field(j)
				value := cv.Field(j).String()
				if len(value) == 0 {
					continue
				}
				switch field.Name {
				case "Name":
					names[value] = true
				case "Email":
					emails[value] = true
				case "Phone":
					phones[value] = true
				case "Postalcode":
					postcodes[value] = true
				case "Organization":
					organizations[value] = true
				}
			}
		}
	}

	// parse into slices
	for k, _ := range emails {
		if k != "" {
			r.Emails = append(r.Emails, k)
		}
	}
	for k, _ := range organizations {
		if k != "" {
			r.Organizations = append(r.Organizations, k)
		}
	}
	for k, _ := range phones {
		if k != "" {
			r.Phones = append(r.Phones, k)
		}
	}
	for k, _ := range postcodes {
		if k != "" {
			r.Postcodes = append(r.Postcodes, k)
		}
	}
	for k, _ := range names {
		if k != "" {
			r.Names = append(r.Names, k)
		}
	}
}
