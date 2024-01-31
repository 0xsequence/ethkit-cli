package main

import (
	"log"

	"github.com/0xsequence/ethkit-cli/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
