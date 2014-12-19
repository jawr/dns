package parser

import (
	"bufio"
	"compress/gzip"
	db "github.com/jawr/dns/database/models/zonefile/parser"
	"github.com/jawr/dns/log"
	"io"
	"os"
)

func (p *Parser) setupFile(filename string, gunzip bool) error {
	_p, err := db.New(filename)
	if err != nil {
		log.Error("Unable to setup Zonefile: %s", err)
		return err
	}
	p.Filename = _p.Filename
	p.TLD = _p.TLD
	p.Date = _p.Date
	p.origin = p.TLD.Name + "."

	err = p.Insert()
	if err != nil {
		log.Error("Unable to save Zonefile Parser: %s", err)
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Error("Unable to open Zonefile: %s", err)
		return err
	}
	p.setupFileDefer = func() {
		file.Close()
	}
	var reader io.Reader = file
	if gunzip {
		reader, err = gzip.NewReader(file)
		if err != nil {
			log.Error("Unable to setup Zonefile gzip reader: %s", err)
			return err
		}
	}
	p.scanner = bufio.NewScanner(reader)
	return nil
}

func (p *Parser) SetupGunzipFile(filename string) error {
	return p.setupFile(filename, true)
}

func (p *Parser) SetupFile(filename string) error {
	return p.setupFile(filename, false)
}
