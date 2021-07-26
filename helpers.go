package main

import "log"

type parameters struct {
	inputPath  string
	outputPath string
	no_header  bool
	compact    bool
	indent     int
}

func processErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
