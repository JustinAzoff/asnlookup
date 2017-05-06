package main

import (
	"fmt"
	"os"

	"github.com/JustinAzoff/asnlookup/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
