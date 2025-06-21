/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/karlrobeck/golang-scraping-template/models"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

// businessCmd represents the business command
var businessCmd = &cobra.Command{
	Use:   "business",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		categoryName := args[0]

		pageStr := cmd.Flag("page").Value.String()

		page, err := strconv.Atoi(pageStr)

		if err != nil {
			log.Fatalf("Invalid page number: %v", err)
		}

		db, err := sql.Open("sqlite", "./local.db")

		if err != nil {
			log.Fatalln(err)
		}

		queries := models.New(db)

		category, err := queries.GetCategory(cmd.Context(), categoryName)

		if err != nil {
			log.Fatalln(err)
		}

		businesses, err := queries.GetBusinessInCategory(cmd.Context(), models.GetBusinessInCategoryParams{
			Name: category.Name,
			Page: int64(page),
		})

		if err != nil {
			log.Fatalln(err)
		}

		if len(businesses) == 0 {
			log.Println("[INFO]:", "No businesses found")

			// check if category is complete before requesting
			if category.IsCompleted == 1 {
				log.Fatalln("No more page avaiable")
			}

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.yellow-pages.ph/%s/page-%d", category.Url, page), nil)

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

			selection := document.Find("h2.search-tradename > a.yp-click")

			if selection.Length() == 0 {
				if _, err := queries.MarkCategoryAsComplete(cmd.Context(), category.ID); err != nil {
					log.Fatalln(err)
				}
				log.Println("No more avaialable businesses")
				return
			}

			for _, node := range selection.EachIter() {

				name := node.Text()

				url, ok := node.Attr("href")

				if !ok {
					continue
				}

				if _, err := queries.CreateBusinessInCategory(cmd.Context(), models.CreateBusinessInCategoryParams{
					CategoryID: category.ID,
					Name:       name,
					Url:        url,
					Page:       int64(page),
				}); err != nil {
					log.Fatalln(err)
				}

			}
			log.Fatalln("[INFO]:", "Fetching complete. please rerun the command again")
		}

		tablePrint := table.NewWriter()

		tablePrint.SetOutputMirror(os.Stdout)

		tablePrint.AppendHeader(table.Row{"#", "Category", "Name", "Url", "Page"})

		for _, business := range businesses {
			tablePrint.AppendRow(table.Row{business.ID, categoryName, business.Name, business.Url, business.Page})
		}

		tablePrint.Render()

	},
}

func init() {
	rootCmd.AddCommand(businessCmd)

	businessCmd.Flags().IntP("page", "p", 1, "Page number to scrape")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// businessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// businessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
