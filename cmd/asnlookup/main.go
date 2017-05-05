package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/JustinAzoff/asnlookup"
)

func main() {

	b, err := asnlookup.NewAsnBackend("asn.db", "asnames.json")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rec, err := b.Lookup(scanner.Text())
		if err != nil {
			log.Print(err)
		} else {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", rec.Prefix, rec.IP, rec.AS, rec.Owner, rec.CC)
		}
	}
}
