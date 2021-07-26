package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

/*
option ideas:
- indent, num spaces
- output to .json file instead of stdout, -o
- take input from file instead of stdin -i
- compact boolean
*/

type parameters struct {
	filepath string
}

type reader interface {
	init(parameters)
	read() ([]string, error)
	cleanup()
}

type fileReader struct {
	file   *os.File
	reader *csv.Reader
}

func (r *fileReader) init(params parameters) {
	file, err := os.Open(params.filepath)
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

func getReader(params parameters) reader {
	var r reader
	if params.filepath == "" {
		r = &stdinReader{}
	} else {
		r = &fileReader{}
	}

	r.init(params)
	return r
}

func processErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func processParams() (parameters, error) {
	// get flag info here ===
	filepath := flag.String("i", "", "Input CSV filepath")

	flag.Parse()

	// get normal arguments here ===

	return parameters{*filepath}, nil
}

func readCsv(params parameters, recordChannel chan<- map[string]string) {
	reader := getReader(params)
	defer reader.cleanup() // cleanup before readCsv() ends

	headers, err := reader.read()
	processErr(err)

	for record, err := reader.read(); err != io.EOF; record, err = reader.read() {
		processErr(err)

		recordMap := make(map[string]string)

		for ind, field := range record {
			recordMap[headers[ind]] = field
		}

		recordChannel <- recordMap
	}

	// tell writeJson() that there are no more records coming
	close(recordChannel)
}

func writeJson(params parameters, recordChannel <-chan map[string]string, done chan<- bool) {
	fmt.Print("[")

	firstRecord := true
	for recordMap, more := <-recordChannel; more; recordMap, more = <-recordChannel {
		var endLastRecord string
		if firstRecord {
			endLastRecord = "\n"
			firstRecord = false
		} else {
			endLastRecord = ",\n"
		}

		fmt.Print(endLastRecord)
		fmt.Print("  {")

		firstField := true
		for header, field := range recordMap {
			var endLastField string
			if firstField {
				endLastField = "\n"
				firstField = false
			} else {
				endLastField = ",\n"
			}

			fmt.Print(endLastField)
			fmt.Printf("    \"%s\": \"%s\"", header, field)
		}

		fmt.Print("\n  }")
	}

	fmt.Println("\n]")

	// tell main() we are done writing the JSON
	done <- true
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	params, err := processParams()
	processErr(err)

	recordChannel := make(chan map[string]string)
	done := make(chan bool)

	go readCsv(params, recordChannel)
	go writeJson(params, recordChannel, done)

	// block main() until writeJson() is done writing
	<-done
}
