package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/tlds"
	"github.com/lib/pq"
	"regexp"
	"strings"
	"time"
)

type Log struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type Parser struct {
	ID       int32       `json:"id"`
	Filename string      `json:"filename"`
	Started  pq.NullTime `json:"started"`
	Finished pq.NullTime `json:"finished"`
	Date     time.Time   `json:"date"`
	TLD      tlds.TLD    `json:"tld"`
	Logs     []Log       `json:"logs"`
}

var tldRE *regexp.Regexp = regexp.MustCompile(`^(\d{8})\-([\w\d\-]+)[\-\.]zone[\-\.](data|gz)`)

const TLD_FILENAME_DATE string = "20060102"

func New(filename string) (Parser, error) {
	args := strings.Split(filename, "/")
	name := args[len(args)-1]
	args = tldRE.FindStringSubmatch(name)
	p := Parser{}
	if len(args) < 4 {
		return p, errors.New("No TLD or date detected in zone filename: " + filename)
	}
	t, err := tlds.GetByName(args[2]).One()
	if err != nil {
		return p, err
	}
	date, err := time.Parse(TLD_FILENAME_DATE, args[1])
	if err != nil {
		return p, err
	}
	p.Filename = name
	p.TLD = t
	p.Date = date
	return p, nil
}

func (p *Parser) Insert() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	var id int32
	err = conn.QueryRow("INSERT INTO parser (filename, parser_date, tld) VALUES ($1, $2, $3) RETURNING id",
		p.Filename,
		p.Date,
		p.TLD.ID,
	).Scan(&id)
	p.ID = id
	return err

}

func (p Parser) Save() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	b, err := json.Marshal(p.Logs)
	if err != nil {
		return err
	}
	_, err = conn.Exec("UPDATE parser SET logs = $1 WHERE id = $2", string(b), p.ID)
	return err
}

func (p *Parser) Update(msg string, args ...interface{}) error {
	log := Log{
		Time:    time.Now(),
		Message: fmt.Sprintf(msg, args...),
	}
	p.Logs = append(p.Logs, log)
	return p.Save()
}

func (p *Parser) Finish() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	t := time.Now()
	p.Finished.Time = t
	_, err = conn.Exec("UPDATE parser SET finished_at = $1 WHERE id = $2", p.Finished, p.ID)
	return err
}
