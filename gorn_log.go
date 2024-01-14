package main

import (
	"fmt"
	"log"
)

// gornLog prints a log with a specific header to stdout depending on logLevel.
//
//	This uses `log.Println` and `log.Fatalln` under the hood.
func gornLog(header LogHeader, v ...any) {
	if header == FATAL && logLevel == FATAL_LEVEL {
		log.Fatal(header, fmt.Sprintln(v...))

	} else if header == WARN && (logLevel == WARN_LEVEL || logLevel == WARN_ONLY) {
		log.Print(header, fmt.Sprintln(v...))

	} else if header == INFO &&  (logLevel == INFO_LEVEL || logLevel == INFO_ONLY) {
		log.Print(header, fmt.Sprintln(v...))

	} else if header == TIME && (logLevel == TIME_LEVEL || logLevel == TIME_ONLY) {
		log.Print(header, fmt.Sprintln(v...))
	}
}
