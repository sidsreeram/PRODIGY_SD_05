package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/gocolly/colly/v2"
)

type Product struct {
    Name   string
    Price  string
    Rating string
}

func main() {
    targetURL := "https://books.toscrape.com/catalogue/category/books_1/index.html"

    file, err := os.Create("products.csv")
    if err != nil {
        log.Fatalf("Could not create file: %v", err)
    }
    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()
    writer.Write([]string{"Product Name", "Price", "Rating"})

    c := colly.NewCollector(
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"),
        colly.Async(true),
        colly.MaxDepth(2),
    )

    c.SetRequestTimeout(180 * time.Second)
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*toscrape.com",
        Parallelism: 1,
        Delay:       5 * time.Second,
    })

    c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
        product := Product{Name: e.ChildText("h3 a")}
        product.Price = e.ChildText(".price_color")
        product.Rating = e.ChildAttr("p.star-rating", "class")
        writer.Write([]string{product.Name, product.Price, product.Rating})
    })

    c.OnError(func(r *colly.Response, err error) {
        log.Printf("Error: %v\nRequest URL: %s", err, r.Request.URL)
    })

    err = c.Visit(targetURL)
    if err != nil {
        log.Fatalf("Failed to scrape: %v", err)
    }

    c.Wait()
    fmt.Println("Scraping complete. Data saved to products.csv")
}
