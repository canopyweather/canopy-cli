/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"canopyweather.com/canopy-cli/internal/config"
	"canopyweather.com/canopy-cli/internal/utils"
	canopyapi "canopyweather.com/canopy-cli/pkg/api"
	"github.com/spf13/cobra"
)

// getSummaryCmd represents the getSummary command
var (
	reportSummaryStartDate string
	reportSummaryEndDate   string
	getSummaryCmd          = &cobra.Command{
		Use:   "getSummary",
		Short: "Creates a csv summary for all states.",
		Long: `Downloads all csvs for a specified date range --startDate and --endDate and merges
them into a single csv.`,
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.GetConfig()
			client := canopyapi.NewClient(conf.ApiKey, conf.Url)

			path, err := os.Getwd()

			if err != nil {
				log.Fatal(err)
			}

			if filepath != "" {
				path = filepath
			}

			finalCSVPath := path + "/" + reportSummaryStartDate + "_" + reportSummaryEndDate + "impact-report-summary.csv"

			if _, err := os.Stat(finalCSVPath); err == nil {
				// File exists, delete it
				os.Remove(finalCSVPath)
			}

			dates, err := utils.GetDatesInRange(reportSummaryStartDate, reportSummaryEndDate)
			if err != nil {
				log.Fatal(err)
			}

			records := make([]canopyapi.ImpactPredictionReport, 0, len(dates))

			for _, date := range dates {
				record, err := client.ImpactPrediction().GetByDate(date)
				if err != nil {
					log.Print(err)
				}
				if record != nil {
					records = append(records, *record)
				}
				time.Sleep(100 * time.Millisecond)
			}

			log.Printf("Total csvs being merged: %d", len(records))

			for index, record := range records {
				stateUrl := record.StateUrl[0]
				if stateUrl == "" {
					continue
				}
				appendSummaryToCSV(stateUrl, finalCSVPath, record.EventDate, index == 0)
			}
		},
	}
)

func init() {
	impactPredictionCmd.AddCommand(getSummaryCmd)

	// Here you will define your flags and configuration settings.
	getSummaryCmd.Flags().StringVarP(&reportSummaryStartDate, "startDate", "s", "", "A start date to pull files for in the format YYYY-MM-DD.")
	getSummaryCmd.Flags().StringVarP(&reportSummaryEndDate, "endDate", "e", "", "An end date to pull files for in the format YYYY-MM-DD.")
	getSummaryCmd.Flags().StringVarP(&filepath, "path", "p", "", "A location to save the file. Defaults to current directory.")
}

func appendSummaryToCSV(url string, outputFilePath string, date string, first bool) error {
	// Step 1: HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Step 2: Read & Parse CSV
	csvReader := csv.NewReader(resp.Body)

	// Step 3: Create Output File
	file, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		if first && record[0] == "State" {
			record = append(record, "Date")
		}

		if !first && record[0] == "State" {
			continue
		}

		if strings.Contains(record[0], "Total") {
			continue
		}

		parsedDate, err := utils.ParseDate(date, "2006-01-02T15:04:05.000Z")
		if err != nil {
			log.Fatal(err)
		}

		if record[0] != "State" {
			record = append(record, utils.FormatDate(parsedDate))
		}

		// Write CSV
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}
	return nil
}
