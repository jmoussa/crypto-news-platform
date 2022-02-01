package scraper

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
	pb "github.com/jmoussa/crypto-dashboard/coindeskmicro/pb"
)

func ScrapeCoinMarketCap() (string, error) {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return "", err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})
	// Instantiate default collector
	c := colly.NewCollector()
	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".currency-name-container"),
			e.ChildText(".col-symbol"),
			e.ChildAttr("a.price", "data-usd"),
			e.ChildAttr("a.volume", "data-usd"),
			e.ChildAttr(".market-cap", "data-usd"),
			e.ChildText(".percent-1h"),
			e.ChildText(".percent-24h"),
			e.ChildText(".percent-7d"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
	return fName, nil
}

func scrapeArticleText(link string) (string, error) {
	scraper := colly.NewCollector()
	res := ""
	scraper.OnResponse(func(r *colly.Response) {
		log.Println("Text Found", r.Request.URL)
	})
	scraper.OnHTML(".at-text p", func(e *colly.HTMLElement) {
		log.Println(e.Text)
		res += "\n" + e.Text
	})
	scraper.Visit(link)
	scraper.Wait()
	return res, nil
}

func ScrapeCoinDeskData(req *pb.GetCoinDeskDataRequest) ([]*pb.Content, error) {
	// Scrape data from CoinDesk
	// max_entries := req.MaxEntries
	items := []*pb.Content{}
	c := colly.NewCollector(
	//colly.AllowedDomains("www.coindesk.com", "coindesk.com"),
	)
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
		//log.Println("Response body:", string(r.Body))
	})
	c.OnHTML("li > a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("Found:", link)
		e.Request.Visit(link)
	})
	c.OnHTML("a[href].headline", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		articleText, err := scrapeArticleText(link)
		if err != nil {
			log.Fatalf("Error when scraping data from %s: %s", link, err)
		}
		// TODO: Send to channel that will stream data to client?
		items = append(items, &pb.Content{
			Title: e.Text,
			Type:  "article",
			Text:  articleText,
			Url:   link,
		})
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Fatalf("Something went wrong: %s", err)
	})
	c.Visit("https://www.coindesk.com/markets")
	c.Wait()
	log.Printf("Scraped %d items", len(items))
	return items, nil
}
