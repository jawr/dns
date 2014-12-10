package parser

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var tldRe *regexp.Regexp = regexp.MustCompile(`^(\d{8})\-([\w\d\-]+)[\-\.]zone[\-\.](data|gz)`)

var tldFilenameDate string = "20060102"

func (p *Parser) setupFile(filename string, gunzip bool) error {
	filenameArgs := strings.Split(filename, "/")
	name := filenameArgs[len(filenameArgs)-1]
	tldNameArgs := tldRe.FindStringSubmatch(name)
	if len(tldNameArgs) < 4 {
		return errors.New("No TLD or date detected in zone filename: " + name)
	}
	p.tldName = tldNameArgs[2]
	var err error
	p.date, err = time.Parse(tldFilenameDate, tldNameArgs[1])
	if err != nil {
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("ERROR: setupFile:Open: %s", err)
		return err
	}
	p.setupFileDefer = func() {
		file.Close()
	}
	var reader io.Reader = file
	if gunzip {
		reader, err = gzip.NewReader(file)
		if err != nil {
			log.Printf("ERROR: setupFile:NewReader: %s", err)
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
