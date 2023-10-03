/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"canopyweather.com/canopy-cli/internal/config"
	canopyapi "canopyweather.com/canopy-cli/pkg/api"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	date   string
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Gets an impact prediction report",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.GetConfig()

			fmt.Printf("CANOPY_API_KEY %s, %s %s \n", conf.ApiKey, "CANOPY_API_URL", conf.Url)

			client := canopyapi.NewClient(conf.ApiKey, conf.Url)

			record, err := client.ImpactPrediction().GetByDate(date)

			if err != nil {
				log.Printf("%s", err)
				return
			}

			if record != nil {
				log.Print("Success")
			}

			path, err := os.Getwd()

			if filepath != "" {
				path = filepath
			}

			if err != nil {
				log.Fatal(err)
			}

			// Download state file.
			stateUrl := record.StateUrl[0]
			if stateUrl != "" {
				stateResp, err := http.Get(record.StateUrl[0])
				if err != nil {
					log.Fatal(err)
				}

				stateFile, err := os.Create(path + "/" + date + "-state-impact-prediction.csv")
				if err != nil {
					log.Fatal(err)
				}
				defer stateFile.Close()

				io.Copy(stateFile, stateResp.Body)

			}

			// Download size file.
			sizeUrl := record.SizeUrl[0]
			if sizeUrl != "" {
				sizeResp, err := http.Get(record.SizeUrl[0])
				if err != nil {
					log.Fatal(err)
				}

				sizeFile, err := os.Create(path + "/" + date + "-size-impact-prediction.csv")
				if err != nil {
					log.Fatal(err)
				}
				defer sizeFile.Close()

				io.Copy(sizeFile, sizeResp.Body)
			}
		},
	}
)

func init() {
	impactPredictionCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.Flags().StringVarP(&date, "date", "d", "", "A date to pull files for in the format YYYY-MM-DD.")
	getCmd.Flags().StringVarP(&filepath, "path", "p", "", "A location to save the files. Defaults to current directory.")
}
