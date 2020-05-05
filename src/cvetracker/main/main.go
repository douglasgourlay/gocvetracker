package main

import (
	"cvetracker/search"
	"cvetracker/updater"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "expected search, updater or write_example")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "updater":
		updater.Main()
		return

	case "search":
		search.Main()
		return

	case "write_example":
		search.WriteExampleConfig("cvesearch.yaml")
		updater.WriteExampleConfig("cveupdater.yaml")
		return

	default:
		fmt.Fprintln(os.Stderr, "expected search or updater")
		os.Exit(1)

	}

}
