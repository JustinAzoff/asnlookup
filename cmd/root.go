package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/JustinAzoff/hostlookup/hostdb"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hostlookup",
	Short: "Hostlookup looks up IP addresses to AS Owner info",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := hostdb.NewHostBackend("shrunken.csv.gz")
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
				fmt.Printf("%s\t%s\n", scanner.Text(), rec.Host)
			}
		}
	},
}
