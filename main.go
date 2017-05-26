package main

import (
	"log"

	"github.com/bcrusu/kcm/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
