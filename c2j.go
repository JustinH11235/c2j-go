package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func processErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// fmt.Println(os.Args)

	// get flag info here ===

	flag.Parse()

	filepath := flag.Arg(0)
	// fmt.Println(filepath)

	file, err := os.Open(filepath)
	processErr(err)

	// close the file before main() ends
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	headers, err := reader.Read()
	processErr(err)

	// fmt.Println(headers)

	// make an empty slice of maps, one map per row
	rowsSlice := []map[string]string{}
	for row, err := reader.Read(); err != io.EOF; row, err = reader.Read() {
		processErr(err)

		// fmt.Println(row)

		rowMap := make(map[string]string)

		for ind, value := range row {
			rowMap[headers[ind]] = value
		}

		rowsSlice = append(rowsSlice, rowMap)
	}

	// fmt.Println(rowsSlice)
	// fmt.Println("\n\n\n\n\n\n\n\n")

	fmt.Println("[")

	for ind, rowMap := range rowsSlice {
		fmt.Println("  {")
		// fmt.Println("   ", rowMap)

		rowMapHeaders := make([]string, 0, len(rowMap))
		for header := range rowMap {
			rowMapHeaders = append(rowMapHeaders, header)
		}

		for ind, header := range rowMapHeaders {
			var endChar string
			if ind < len(rowMap)-1 {
				endChar = ","
			} else {
				endChar = ""
			}

			fmt.Printf("    \"%s\": \"%s\"%s\n", header, rowMap[header], endChar)
		}

		var endChar string
		if ind < len(rowsSlice)-1 {
			endChar = ","
		} else {
			endChar = ""
		}

		fmt.Printf("  }%s\n", endChar)
	}

	fmt.Println("]")
}
