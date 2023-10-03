/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"canopyweather.com/canopy-cli/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var filepath string
var rootCmd = &cobra.Command{
	Use:   "canopy-cli",
	Short: "An application for interacting with the canopy apis.",
	Long: `The canopy-cli allows users to interact with the canopy api using their api key. 
	
It also provides common functionality on top of those apis for specific use-cases that 
require the joining of multiple api requests.
	
Outputs can vary depending on the command, but a majority of the time the system 
will write a file output.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	config.Setup()
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.canopy-cli.yaml)")
	getCmd.PersistentFlags().StringVarP(&filepath, "path", "p", "", "A location to save the files. Defaults to current directory.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
