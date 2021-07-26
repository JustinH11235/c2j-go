package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func readCsv(params parameters, recordChannel chan<- map[string]string) {
	reader := createReader(params)
	defer reader.cleanup() // cleanup before readCsv() ends

	// if csv contains headers, read headers
	var headers []string
	var err error
	if !params.no_header {
		headers, err = reader.read()
		processErr(err)
	}

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
		jsonBytes, err := json.MarshalIndent(recordMap, indent, indent) // converts map into JSON string
		processErr(err)
		return indent + string(jsonBytes) // add indent to first line
	}
}

func writeJson(params parameters, recordChannel <-chan map[string]string, done chan<- bool) {
	writer := createWriter(params)
	defer writer.cleanup() // cleanup before writeJson() ends

	var lineBreak string
	if params.compact {
		lineBreak = ""
	} else {
		lineBreak = "\n"
	}

	writer.write("[")

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
		writer.write(endLastRecord)

		writer.write(mapToJsonString(recordMap, params))

		/*
			Old code that implemented json.Marshal(), removed since json.Marshal() handles all the quirks
			of JSON such as escaping
		*/

		// writer.write("  {")

		// firstField := true
		// for header, field := range recordMap {
		// 	var endLastField string
		// 	if firstField {
		// 		endLastField = "\n"
		// 		firstField = false
		// 	} else {
		// 		endLastField = ",\n"
		// 	}

		// 	writer.write(endLastField)
		//  //writer.write(fmt.Sprintf("    \"%s\": \"%s\"", header, field))
		// 	// fmt.Printf("    \"%s\": \"%s\"", header, field)
		// }

		// writer.write("\n  }")
	}

	writer.write(lineBreak + "]")

	// tell main() we are done writing the JSON
	done <- true
}

func processParams() (parameters, error) {
	// get flag info here ===
	filepath := flag.String("i", "", "Input CSV filepath")
	output := flag.String("o", "", "Output JSON filepath (WARNING: this overwrites files with same name!)")
	no_header := flag.Bool("no-header", false, "Input does not contain header (integers will be used as headers)")
	compact := flag.Bool("compact", false, "Output without whitespace")
	indent := flag.Int("indent", 2, "Indent of lines (always 0 if -compact is set)")

	flag.Parse()

	// get normal arguments here ===

	return parameters{*filepath, *output, *no_header, *compact, *indent}, nil
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

	// readCsv() as one goroutine and send data to write to writeJson() through recordChannel
	go readCsv(params, recordChannel)
	go writeJson(params, recordChannel, done)

	// block main() until writeJson() is done writing
	<-done
}
