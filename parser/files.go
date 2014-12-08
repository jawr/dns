package parser

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func (p *Parser) setupFile(filename string, gunzip bool) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("ERROR: setupFile:Open: %s", err)
		return err
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
