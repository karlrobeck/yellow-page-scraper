/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/karlrobeck/golang-scraping-template/models"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		canonicalLink := args[0]

		db, err := sql.Open("sqlite", "./local.db")

		if err != nil {
			log.Fatalln(err)
		}

		queries := models.New(db)

		info, err := queries.GetBusinessInfo(cmd.Context(), canonicalLink)

		if err != nil {
			log.Println(err)

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.yellow-pages.ph/%s", canonicalLink), nil)

			if err != nil {
				log.Fatalln(err)
			}

			// Add headers
			request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0")
			request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			request.Header.Add("Accept-Language", "en-US,en;q=0.5")
			request.Header.Add("Upgrade-Insecure-Requests", "1")
			request.Header.Add("Sec-Fetch-Dest", "document")
			request.Header.Add("Sec-Fetch-Mode", "navigate")
			request.Header.Add("Sec-Fetch-Site", "cross-site")
			request.Header.Add("Priority", "u=0, i")
			request.Header.Add("Pragma", "no-cache")
			request.Header.Add("Cache-Control", "no-cache")

			// Set referrer
			request.Header.Add("Referer", "https://www.google.com/")

			response, err := http.DefaultClient.Do(request)

			if err != nil {
				log.Fatalln(err)
			}

			if response.StatusCode != 200 {
				log.Fatalln(response)
			}

			document, err := goquery.NewDocumentFromReader(response.Body)

			if err != nil {
				log.Fatalln(err)
			}

			tradeName := document.Find("h1.h1-tradename")
			businessName := document.Find("h2.h2-businessname")
			addresses := document.Find(`a[data-section="address"]`).Map(func(i int, s *goquery.Selection) string {
				return s.Text()
			})
			phoneNumbers := document.Find(`a[href^="tel:"]`).Map(func(i int, s *goquery.Selection) string {
				return s.AttrOr("href", "")
			})
			emails := document.Find(`a[href^="mailto:"]`).Map(func(i int, s *goquery.Selection) string {
				return s.AttrOr("href", "")
			})
			websites := document.Find(`a[data-section="website"]`).Map(func(i int, s *goquery.Selection) string {
				return s.AttrOr("href", "")
			})
			socials := document.Find(`a[data-section="social"]`).Map(func(i int, s *goquery.Selection) string {
				return s.AttrOr("href", "")
			})
			ratingStr := document.Find(".rating-num").Text()
			var rating float64

			if ratingStr != "" {
				if r, err := strconv.ParseFloat(strings.TrimSpace(ratingStr), 32); err == nil {
					rating = r
				} else {
					log.Printf("Failed to parse rating: %v", err)
				}
			}

			if _, err := queries.CreateBusinessInfo(cmd.Context(), models.CreateBusinessInfoParams{
				TradeName:     sql.NullString{String: tradeName.Text(), Valid: true},
				BusinessName:  sql.NullString{String: businessName.Text(), Valid: true},
				Address:       sql.NullString{String: strings.Join(addresses, ","), Valid: true},
				PhoneNumber:   sql.NullString{String: strings.Join(phoneNumbers, ","), Valid: true},
				Email:         sql.NullString{String: strings.Join(emails, ","), Valid: true},
				Website:       sql.NullString{String: strings.Join(websites, ", "), Valid: true},
				SocialMedia:   sql.NullString{String: strings.Join(socials, ","), Valid: true},
				CanonicalLink: canonicalLink,
				Rating:        sql.NullFloat64{Float64: rating},
			}); err != nil {
				log.Fatalln(err)
			}

			fmt.Println(tradeName.Text(), businessName.Text(), phoneNumbers, emails, websites, socials, addresses, rating)

		}

		fmt.Println(info)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
