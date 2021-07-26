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

## Credits
Based on [this](https://levelup.gitconnected.com/tutorial-how-to-create-a-cli-tool-in-golang-a0fd980264f) tutorial.


