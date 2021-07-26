package main

import (
	"fmt"
	"os"
)

type writer interface {
	init(parameters)
	write(string)
	cleanup()
}

type stdinWriter struct{}

func (w *stdinWriter) init(params parameters) {}

func (w *stdinWriter) write(out string) {
	_, err := fmt.Print(out)
	processErr(err)
}

func (w *stdinWriter) cleanup() {}

type fileWriter struct {
	file *os.File
}

func (w *fileWriter) init(params parameters) {
	file, err := os.Create(params.outputPath)
	processErr(err)

	w.file = file
}

func (w *fileWriter) write(out string) {
	_, err := w.file.WriteString(out)
	processErr(err)
}

func (w *fileWriter) cleanup() {
	w.file.Close()
}

func createWriter(params parameters) writer {
	var w writer
	if params.outputPath == "" {
		w = &stdinWriter{}
	} else {
		w = &fileWriter{}
	}

	w.init(params)
	return w
}
