package targets

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosecondhand/src/database"
	"gosecondhand/src/utils"
	"log"
	"strconv"
	"strings"
	"sync"
)

// General algorithm of this web-scraper in the greater system.
// 1. Get HTML code of the search page
// 2. Parse through and find data
// 3. Get HTML code of each item, using the found URL and find more data there
// 4. Store all the data in the database.
// NOTE: Currently only searches through the first search result page.

// GenerateEbayURL returns a valid url to an Ebay search result page given a search input.
// For instance, calling the function with "Beatles puzzle" will return the url to the
// search results when searching "Bealtes puzzle" on Ebay.
// Example:
//
//    searchpageURL := GenerateEbayURL("Beatles puzzle")

func GenerateEbayURL(searchString string) string {
	concatenatedString := ("https://www.ebay.com/sch/i.html?_from=R40&_nkw=" + searchString + "&sacat=0&_ipg=200")
	return strings.ReplaceAll(concatenatedString, " ", "+")
}

// EbayTrimUrl will remove unneccesary text at the end of a URL
// to an Ebay item.
// Example:
//
// trimmedURL := EbayTrimUrl("https://www.ebay.com/itm/233732513763?epid=1305260978&_trkparms=ispr%3D1&hash=item366b8b1fe3:g:nGUAAOSw6a5fdvPB&amdata=enc%3AAQAGAAAA4NtZmuG9imn43XDXjwNuAkmo7M%2F5WpbaU2wu9WnhxW7vexMR%2B79O1cwW0vAUCUQNuhwYTwpC4GU6BtbLh61dqB6AY62YoCf2YYwg4EBkVCxB%2FlS348heP4SXhddV4%2Bt78MNMtBMHp11TZqgHK50oeZFJdQt0wpe0I8khQt9GNPeJBEIRFHjImdizNEp6qutvqZRFAgjgV5fpKNztXJNC3LsU86Br2ONqgq5q8Z%2BBxuXK6yehcL1VFX0xrQCnZk")
//
// (trimmedURL ==  https://www.ebay.com/itm/233732513763)
func EbayTrimUrl(url string) string {
	result := ""
	runeUrl := []rune(url)

	for i := 0; i < len(runeUrl); i++ {
		if string(runeUrl[i]) == "?" {
			break
		}
		result += string(runeUrl[i])
	}
	return result
}

// trimEbayCategory will take a string and return it without newline (\n) and anything after it.
// Example:
//
// category := trimEbayCategory("Video Games\n
//                               Playstation Games")
func trimEbayCategory(categories string) string {
	result := ""
	runeCategories := []rune(categories)

	for i := 0; i < len(runeCategories); i++ {
		currentChar := string(runeCategories[i])
		if currentChar == "\n" {
			break
		}

		result += currentChar
	}
	return result
}

// trimEbayPrice will return a string but without commas and anything after a space.
// Example:
//
// price := trimEbayPrice("1,200 USD")
// (price == 1200)
func trimEbayPrice(price string) string {
	result := ""
	runePrice := []rune(price)

	for i := 0; i < len(runePrice); i++ {
		currentChar := string(runePrice[i])

		if currentChar == " " {
			return result
		} else if currentChar == "," {
			result += ""
		} else {
			result += currentChar
		}

	}
	return result
}

// trimDescription will return a string but without the word "Read" anything after it.
// Here, it is used to shorten the generic condition description Ebay has for items; new, used and so on.
// Example:
//
// desc := trimDescription("New: A brand-new, unused, unopened, undamaged item (including handmade items). See the seller's ...
//                          Read more")
func trimDescription(desc string) string {
	split := strings.Split(desc, "Read")
	if len(split) > 1 {
		return split[0]
	} else {
		return desc
	}
}

// ScrapeEbayPageData will, given a goquery-document and a string
// containing a search word, scrape the Ebay search page for the
// word and scrape each item page on it.
// After scraping the pages for data, it will be stored in a MySQL-database.
//
// The data which will be scraped is each item's:
// URL, Picture-URL, Title, Description, Price and Category.
//
// All of the above will be stored in the database along with the string
// argument and an "Ebay"-string.
//
// Example:
//
// document, error := goquery.NewDocumentFromReader(response.Body)
// ...
// ScrapeEbayPageData(document, "Beatles Puzzle")
func ScrapeEbayPageData(document *goquery.Document, searchString string) []utils.Item {
	// ALLOCATE ITEM STRUCT
	data := utils.Item{}

	// Variable for storing all items as a list
	var dataList []utils.Item

	// DATABASE CONNECTION
	db := database.ConnectDB()

	// Suspend execution of this function until surround has ran
	defer db.Close()

	// ERROR HANDLING FOR DATABSE CONNECTION
	err := db.Ping()
	if err != nil {
		log.Fatal("Connection could not be verified with Ping(): ", err)
	}

	table_name := "Ebay_" + strings.ReplaceAll(searchString, " ", "_")
	err = database.CreateTables(db, table_name)
	// ERROR HANDLING FOR TABLE CREATION FOR DATABASE
	if err != nil {
		log.Fatal("Could not create tables: ", err)
	}

	document.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		found := item.Find("a.s-item__link")
		title := found.Text()
		url, _ := found.Attr("href")
		processedUrl := EbayTrimUrl(url)

		price_span := item.Find("span.s-item__price").Text()
		processedPrice := trimEbayPrice(price_span)

		price, _ := strconv.ParseFloat(processedPrice, 64)

		// Insert item metdata into struct object
		data.SearchString = searchString
		data.Site = "Ebay"
		data.URL = processedUrl
		data.Title = title
		data.Price = price

		data = ScrapeEbayItemPageData(processedUrl, data)

		// Append item to list of all items
		dataList = append(dataList, data)
		// Insert to database
		database.InsertData(db, table_name, data)
	})
	// Return the list of items scraped from Ebay using item struct from utils.go
	return dataList
}

// ScrapeEbayItemPage will scrape an Ebay item page given its URL as a string,
// and an Item-struct called 'data' for storing the scraped data in.
// The function will return 'data' but with PictureURL, Category and
// Description added.
//
// Example:
//
// data := utils.Item{}
// url := "https://www.ebay.com/itm/233732513763"
// data = ScrapeItemPageData(url, data)
func ScrapeEbayItemPageData(url string, data utils.Item) utils.Item {
	response := utils.GetHTML(url)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	categories := document.Find("li#vi-VR-brumb-lnkLst")
	category := trimEbayCategory(strings.TrimSpace(categories.Find("a.thrd").Text()))

	description := document.Find("span.ux-expandable-textual-display-block-inline").Text()
	description = trimDescription(description)

	//Extract src
	picture := document.Find("img#icImg")
	pictureURL, _ := picture.Attr("src")

	data.PictureURL = pictureURL
	data.Category = category
	data.Description = description

	return data
}

// EbayResponse will given a search word as a string,
// get a HTML response from the search page and scrape its items' data.
// The data will then be stored in a MySQL-database.
// Example:
//
// EbayResponse("Beatles Puzzle")
func EbayResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()
	ebayURL := GenerateEbayURL(searchString)

	response := utils.GetHTML(ebayURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	ScrapeEbayPageData(document, searchString)
	fmt.Println("EBAY RESPONSE COMPLETE!")
}
