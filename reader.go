package main

import (
	"encoding/csv"
	"os"
)

type reader interface {
	init(parameters)
	read() ([]string, error)
	cleanup()
}

type stdinReader struct {
	reader *csv.Reader
}

func (r *stdinReader) init(params parameters) {
	reader := csv.NewReader(os.Stdin)
	reader.Comma = ','
	// reader.LazyQuotes = true

	r.reader = reader
}

func (r *stdinReader) read() ([]string, error) {
	return r.reader.Read()
}

func (r *stdinReader) cleanup() {}

type fileReader struct {
	file   *os.File
	reader *csv.Reader
}

func (r *fileReader) init(params parameters) {
	// check if file exists
	_, err := os.Stat(params.inputPath)
	processErr(err)

	file, err := os.Open(params.inputPath)
	processErr(err)

	reader := csv.NewReader(file)
	reader.Comma = ','
	// reader.LazyQuotes = true

	r.file = file
	r.reader = reader
}

func (r *fileReader) read() ([]string, error) {
	return r.reader.Read()
}

func (r *fileReader) cleanup() {
	r.file.Close()
}

func createReader(params parameters) reader {
	var r reader
	if params.inputPath == "" {
		r = &stdinReader{}
	} else {
		r = &fileReader{}
	}

	r.init(params)
	return r
}
