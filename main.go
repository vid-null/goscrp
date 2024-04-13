package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for _, letter := range alphabet {
		url := fmt.Sprintf("https://www.reginamaria.ro/utile/dictionar-de-analize/%c?litera=%c", letter, letter)
		scrape(url)
	}
}

func scrape(url string) {
	c := colly.NewCollector()

	var pages []string

	c.OnHTML(".views-row .views-field-title a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absURL := e.Request.AbsoluteURL(link) // Convert relative URL to absolute URL
		pages = append(pages, absURL)
	})

	downloadPage := func(url string) {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error downloading page:", err)
			return
		}
		defer resp.Body.Close()

		downloadDir := "./download"
		err = os.MkdirAll(downloadDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		filename := filepath.Join(downloadDir, filepath.Base(url))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		c := colly.NewCollector()
		//c.OnHTML(".region-content", func(e *colly.HTMLElement) {
		//	// Convert the selected HTML element to string
		//	var sb strings.Builder
		//	err := html.Render(&sb, e.DOM.Nodes[0])
		//	if err != nil {
		//		fmt.Println("Error rendering HTML:", err)
		//		return
		//	}
		//
		//	// Write the content to the file
		//	_, err = file.WriteString(sb.String())
		//	if err != nil {
		//		fmt.Println("Error writing to file:", err)
		//		return
		//	}
		//	fmt.Println("HTML content extracted and saved from:", url)
		//})

		c.OnHTML(".region-content", func(e *colly.HTMLElement) {
			// Write the text content to the file
			_, err := file.WriteString(e.Text)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Println("Text content extracted and saved from:", url)
		})

		// Visit the page to extract content
		err = c.Visit(url)
		if err != nil {
			fmt.Println("Error visiting page:", err)
			return
		}
	}

	// Visit the website and collect links
	c.Visit(url)

	// Download each page found on the specified path
	for _, pageURL := range pages {
		downloadPage(pageURL)
	}
}
