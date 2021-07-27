# c2j-go

## About
A CSV to JSON converter command line tool written in [Go](https://golang.org/).

## Highlights
* Utilizes Go's built-in concurrency support through [Goroutines](https://gobyexample.com/goroutines) and [channels](https://gobyexample.com/channels)!
* For abstracting reading and writing (to support using stdin/stdout v.s. files), used an [interface](https://gobyexample.com/interfaces) approach.

## Options/flags supported
* -i \<filepath\>:\
        Input CSV filepath
* -o \<filepath\>:\
        Output JSON filepath (WARNING: this overwrites files with same name!)
* --no-header:\
        Input does not contain header (integers will be used as headers)
* --compact:\
        Output without whitespace
* --indent \<spaces\>:\
        Indent of lines (always 0 if -compact is set) (default 2)

## Installation Instructions
1. Make sure Go is installed, instructions [here](https://golang.org/doc/install) for your OS
2. Clone this repository and navigate to it in the terminal
3. Run `go build`
4. Use the newly compiled binary: `./c2j-go -i test.csv`
5. (optionally put this binary into a directory that your system's PATH variable looks in so that you can run it from anywhere using `c2j-go -i test.csv`. For Linux you can do `sudo cp ./c2j-go /usr/local/bin/`)

## Credits
Based on [this](https://levelup.gitconnected.com/tutorial-how-to-create-a-cli-tool-in-golang-a0fd980264f) tutorial.


