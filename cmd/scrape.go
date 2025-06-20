/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/karlrobeck/golang-scraping-template/internal"
	"github.com/spf13/cobra"
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// config
		content, err := os.ReadFile("./config.toml")

		if err != nil {
			log.Fatalln(err)
		}

		var config internal.SystemConfig

		if _, err = toml.Decode(string(content), &config); err != nil {
			log.Fatalln(err)
		}

		// scrape the data

		// save to database

	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
