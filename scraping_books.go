package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type books_info struct {
	title, title_link, image, price string
}

func main() {
	c := colly.NewCollector()

	// initializing the slice of structs to store the data to scrape
	var books []books_info

	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		book := books_info{}
		book.title = h.ChildAttr("h3 a", "title")
		book.title_link = h.ChildAttr("h3 a", "href")
		book.image = h.ChildAttr("div.image_container a img", "src")
		book.price = h.ChildText("p.price_color")
		books = append(books, book)
	})

	c.Visit("https://books.toscrape.com/")

	// --------------------------------------------------------------------------
	file, err := os.Create("books.csv")
	if err != nil {
		log.Fatalln("Failed to create an output csv file.", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"title",
		"title-link",
		"image",
		"price",
	}

	writer.Write(headers)

	for _, book := range books {
		record := []string{
			book.title,
			book.title_link,
			book.image,
			book.price,
		}
		writer.Write(record)
	}

	defer writer.Flush()
}
