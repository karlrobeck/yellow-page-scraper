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

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/karlrobeck/golang-scraping-template/models"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

// categoriesCmd represents the categories command
var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		db, err := sql.Open("sqlite", "./local.db")

		if err != nil {
			log.Fatalln(err)
		}

		queries := models.New(db)

		categories, err := queries.GetAllCategories(cmd.Context())

		if err != nil {
			log.Fatalln(err)
		}

		// request if categories are empty
		if len(categories) == 0 {

			log.Println("[INFO]:", "Empty categories.", "requesting...")

			request, err := http.NewRequest(http.MethodGet, "https://www.yellow-pages.ph/category", nil)
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

			if response.StatusCode == 200 {

				document, err := goquery.NewDocumentFromReader(response.Body)

				if err != nil {
					log.Fatalln(err)
				}
				for _, node := range document.Find("a.category-item").EachIter() {

					link, ok := node.Attr("href")

					if !ok {
						continue
					}

					name, ok := node.Attr("data-type")

					if !ok {
						continue
					}

					sizeStr := node.Children().Last().Text()

					var size int
					fmt.Sscanf(sizeStr, "%d", &size)

					// save to database
					if _, err := queries.CreateCategory(cmd.Context(), models.CreateCategoryParams{
						Name: name,
						Url:  link,
						Size: int64(size),
					}); err != nil {
						fmt.Println(err)
					}

				}

				log.Println("[INFO]:", "Request complete. rerun again to see the results")

			} else {
				log.Fatalln("[INFO]:", "Request failed")
			}

		} else {

			tablePrint := table.NewWriter()

			tablePrint.SetOutputMirror(os.Stdout)

			tablePrint.AppendHeader(table.Row{"#", "Name", "URL", "Size", "Is Complete"})

			for _, category := range categories {
				tablePrint.AppendRow(table.Row{category.ID, category.Name, category.Url, category.Size, category.IsCompleted})
			}

			tablePrint.Render()

		}

	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categoriesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categoriesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
