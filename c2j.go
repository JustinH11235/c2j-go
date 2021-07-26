package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/*
option ideas:
- output to .json file instead of stdout, -o
- header boolean
*/

type parameters struct {
	filepath  string
	no_header bool
	compact   bool
	indent    int
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

func readCsv(params parameters, recordChannel chan<- map[string]string) {
	reader := getReader(params)
	defer reader.cleanup() // cleanup before readCsv() ends

	headers, err := reader.read()
	processErr(err)

	for record, err := reader.read(); err != io.EOF; record, err = reader.read() {
		processErr(err)

		recordMap := make(map[string]string)

		for ind, field := range record {
			if params.no_header {
				recordMap[fmt.Sprint(ind)] = field
			} else {
				recordMap[headers[ind]] = field
			}
		}

		recordChannel <- recordMap
	}

	// tell writeJson() that there are no more records coming
	close(recordChannel)
}

func mapToJsonString(recordMap map[string]string, params parameters) string {
	indent := strings.Repeat(" ", params.indent)

	if params.compact {
		jsonBytes, err := json.Marshal(recordMap)
		processErr(err)
		return string(jsonBytes)
	} else {
		jsonBytes, err := json.MarshalIndent(recordMap, indent, indent)
		processErr(err)
		return indent + string(jsonBytes) // add indent to first line
	}
}

func writeJson(params parameters, recordChannel <-chan map[string]string, done chan<- bool) {
	var lineBreak string
	if params.compact {
		lineBreak = ""
	} else {
		lineBreak = "\n"
	}

	fmt.Print("[")

	firstRecord := true
	for recordMap, more := <-recordChannel; more; recordMap, more = <-recordChannel {
		var lineEnd string
		if firstRecord {
			lineEnd = ""
			firstRecord = false
		} else {
			lineEnd = ","
		}

		endLastRecord := lineEnd + lineBreak
		fmt.Print(endLastRecord)

		// convert map into JSON string
		fmt.Print(mapToJsonString(recordMap, params))

		// fmt.Print("  {")

		// firstField := true
		// for header, field := range recordMap {
		// 	var endLastField string
		// 	if firstField {
		// 		endLastField = "\n"
		// 		firstField = false
		// 	} else {
		// 		endLastField = ",\n"
		// 	}

		// 	fmt.Print(endLastField)
		// 	// fmt.Printf("    \"%s\": \"%s\"", header, field)
		// }

		// fmt.Print("\n  }")
	}

	fmt.Println(lineBreak + "]")

	// tell main() we are done writing the JSON
	done <- true
}

func processParams() (parameters, error) {
	// get flag info here ===
	filepath := flag.String("i", "", "Input CSV filepath")
	no_header := flag.Bool("no-header", false, "Input does not contain header (integers will be used as headers)")
	compact := flag.Bool("compact", false, "Output without whitespace")
	indent := flag.Int("indent", 2, "Indent of lines (always 0 if -compact is set)")

	flag.Parse()

	// get normal arguments here ===

	return parameters{*filepath, *no_header, *compact, *indent}, nil
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
