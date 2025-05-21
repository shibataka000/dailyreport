package main

import (
	"log"

	"github.com/shibataka000/dailyreport/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
