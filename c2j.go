package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type parameters struct {
	filepath string
}

func processErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func processParams() (parameters, error) {
	// fmt.Println(os.Args)

	// get flag info here ===

	flag.Parse()

	filepath := flag.Arg(0) // change to accept stdin by default, option for filepath in future
	// fmt.Println(filepath)

	return parameters{filepath}, nil
}

func makeReader(params parameters) (*csv.Reader, *os.File) { // later add a method that reads differently based on input location, also have a close method
	file, err := os.Open(params.filepath)
	processErr(err)

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	return reader, file
}

func readCsv(params parameters, rowChannel chan<- map[string]string) {
	reader, file := makeReader(params)
	defer file.Close() // close the file before readCsv() ends

	headers, err := reader.Read()
	processErr(err)

	for row, err := reader.Read(); err != io.EOF; row, err = reader.Read() {
		processErr(err)

		rowMap := make(map[string]string)

		for ind, field := range row {
			rowMap[headers[ind]] = field
		}

		rowChannel <- rowMap
	}

	// tell writeJson() that there are no more rows coming
	close(rowChannel)
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
		fmt.Printf("Usage: %s <filepath>\n", os.Args[0])
		flag.PrintDefaults()
	}

	params, err := processParams()
	processErr(err)

	rowChannel := make(chan map[string]string)
	done := make(chan bool)

	go readCsv(params, rowChannel)
	go writeJson(params, rowChannel, done)

	<-done
}
