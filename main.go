/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"canopyweather.com/canopy-cli/cmd"
)

func main() {
	if os.Getenv("CANOPY_API_KEY") == "" {
		log.Fatal("CANOPY_API_KEY must be set to a valid api key")
	}

	cmd.Execute()
}
