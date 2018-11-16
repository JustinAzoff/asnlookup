package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/JustinAzoff/asnlookup/asndb"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "asnlookup",
	Short: "Asnlookup looks up IP addresses to AS Owner info",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := asndb.NewAsnBackend("asn.db", "asnames.json")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Reading IP addresses from stdin...")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			rec, err := b.Lookup(scanner.Text())
			if err != nil {
				log.Print(err)
			} else {
				fmt.Printf("%s\t%s\t%d\t%s\t%s\n", rec.Prefix, rec.IP, rec.AS, rec.CC, rec.Owner)
			}
		}
	},
}
