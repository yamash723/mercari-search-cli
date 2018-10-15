package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"regexp"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Item is a summarize information for mercari.jp item
type Item struct {
	Name      string
	PageURL   string
	ImageURL  string
	Price     int
	OnSale    bool
}
// SearchCondition is a conditions for mercari.jp
type SearchCondition struct {
	Keyword       string
	BrandID       uint
	BrandName     string
	CategoryRoot  uint
	CategoryChild uint
	PriceMin      uint
	PriceMax      uint
	Page          uint
	ItemCondition string
	ShippingPayer uint
	OnSale        bool
	SortByDesc    bool
}

func main() {
	app := cli.NewApp()
	app.Name = "mercari-search"
	app.Usage = "Fetch a search result from mercari.jp"
	app.Action = execute
	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "page, p",
			Usage: "page number",
		},

		cli.StringFlag{
			Name:  "keyword, k",
			Usage: "search keyword",
		},

		cli.UintFlag{
			Name:  "price-min, min",
			Usage: "minimum price",
		},

		cli.UintFlag{
			Name:  "price-max, max",
			Usage: "maximum price",
		},

		cli.UintFlag{
			Name:  "category-root",
			Usage: "category root number",
		},

		cli.UintFlag{
			Name:  "category-child",
			Usage: "category child number",
		},

		cli.StringFlag{
			Name:  "brand-name",
			Usage: "brand name keyword",
		},

		cli.UintFlag{
			Name:  "brand-id",
			Usage: "brand id number",
		},

		cli.BoolFlag{
			Name:  "desc",
			Usage: "search in desc order",
		},

		cli.BoolFlag{
			Name:  "on-sale",
			Usage: "fetch only a on-sale items",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(c *cli.Context) error {
	condition := SearchCondition{
		Keyword: c.String("keyword"),
		BrandName: c.String("brand-name"),
		SortByDesc: c.Bool("desc"),
		Page: c.Uint("page"),
		CategoryRoot: c.Uint("category-root"),
		CategoryChild: c.Uint("category-child"),
		BrandID: c.Uint("brand-id"),
		PriceMin: c.Uint("price-min"),
		PriceMax: c.Uint("price-max"),
		OnSale: c.Bool("on-sale"),
	}

	items, err := fetchSearchResult(condition)
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Println(strings.Repeat("-", 100))
		fmt.Println("Name:      ", item.Name)
		fmt.Println("Price:     ", item.Price)
		fmt.Println("OnSale:    ", item.OnSale)
		fmt.Println("PageURL:   ", item.PageURL)
		fmt.Println("ImageURL:  ", item.ImageURL)
	}

	return nil
}

func fetchSearchResult(condition SearchCondition) ([]Item, error) {
	url := buildURL(condition)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed fetch a search result from mercari.jp")
	}

	items := make([]Item, 0)

	headline := doc.Find(".search-result-head").First().Text()
	if strings.Contains(headline, "検索結果 0件") {
		return items, nil
	}

	doc.Find(".items-box").Each(func(_ int, s *goquery.Selection) {
		name := s.Find(".items-box-name").First().Text()
		pageURL, _ := s.Find("a").First().Attr("href")
		imageURL, _ := s.Find(".items-box-photo > img").First().Attr("data-src")
		OnSale := s.Find(".item-sold-out-badge").First().Text() == ""

		priceText := s.Find(".items-box-price").First().Text()
		priceText = regexp.MustCompile(`[^\d]`).ReplaceAllString(priceText, "")
		price, _ := strconv.Atoi(priceText)

		item := Item{
			Name:      name,
			Price:     price,
			PageURL:   pageURL,
			ImageURL:  imageURL,
			OnSale:    OnSale,
		}

		items = append(items, item)
	})

	return items, nil
}

func buildURL(condition SearchCondition) string {
	params := map[string]string{}
	params["page"] = uintToStringForParams(condition.Page)
	params["keyword"] = url.QueryEscape(condition.Keyword)
	params["category-root"] = uintToStringForParams(condition.CategoryRoot)
	params["category-child"] = uintToStringForParams(condition.CategoryChild)
	params["brand-name"] = url.QueryEscape(condition.BrandName)
	params["brand-id"] = uintToStringForParams(condition.BrandID)
	params["price-min"] = uintToStringForParams(condition.PriceMin)
	params["price-max"] = uintToStringForParams(condition.PriceMax)

	if condition.SortByDesc {
		params["sort_order"] = "created_desc"
	} else {
		params["sort_order"] = "created_asc"
	}

	if condition.OnSale {
		params["status_on_sale"] = "1"
	}

	buf := make([]byte, 0)
	for key, value := range params {
		if len(buf) != 0 {
			buf = append(buf, "&"...)
		}

		buf = append(buf, (key + "=" + value)...)
	}

	return "https://www.mercari.com/jp/search/" + "?" + string(buf)
}

func uintToStringForParams(value uint) string {
	if value == 0 {
		return ""
	}

	return strconv.Itoa(int(value))
}
