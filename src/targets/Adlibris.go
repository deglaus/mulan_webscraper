package targets

import (
	"fmt"
	"gosecondhand/src/database"
	"gosecondhand/src/utils"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// GenerateAdlibrisURL returns a valid url to a Adlibris search result page given a search input.
// For instance, calling the function with "berserk manga" will return the url to the
// search results when searching "berserk manga" on Adlibris.
// Example:
//
//    searchpageURL := GenerateEbayURL("berserk manga")

func GenerateAdlibrisURL(searchString string) string {
	concatenatedString := ("https://www.adlibris.com/se/sok?q=" + searchString)
	return strings.ReplaceAll(concatenatedString, " ", "%20")
}

func ScrapeAdlibrisPageData(document *goquery.Document, searchString string) []utils.Item {
	// ALLOCATE A NEW DATA ITEM OBJECT
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

	table_name := "Adlibris_" + strings.ReplaceAll(searchString, " ", "_")

	document.Find("div.search-result__list-view__product").Each(func(index int, item *goquery.Selection) {
		found := item.Find("a.search-result__product__name")
		title := found.Text()
		processedTitle := strings.TrimSpace(title)
		url, _ := found.Attr("href")
		processedUrl := "https://www.adlibris.com" + url

		price_div := item.Find("div.price")
		price := strings.TrimSpace(price_div.Text())
		price = strings.ReplaceAll(price, " kr", "")
		convertedPrice, _ := strconv.ParseFloat(price, 32)

		pic_div := item.Find("div.search-result__list-view__product__image-and-information-container")
		pic := pic_div.Find("img")
		picUrl, _ := pic.Attr("data-src")

		desc_div := item.Find("p.search-result__list-view__product__information__description")
		desc := desc_div.Text()

		// INSERTS APPROPRIATE DATA TO STRUCT FIELDS
		data.SearchString = searchString
		data.Site = "Adlibris"
		data.URL = processedUrl
		data.PictureURL = picUrl
		data.Title = processedTitle
		data.Description = desc
		data.Price = convertedPrice
		// Append item to list of all items
		dataList = append(dataList, data)

		// INSERT DATA TO DATABASE
		database.InsertData(db, table_name, data)
	})
	return dataList
}

// AdlibrisResponse will given a search word as a string,
// get a HTML response from the search page and scrape its items' data
// Example:
//
// AdlibrisResponse("berserk manga")

func AdlibrisResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()
	adlibrisURL := GenerateAdlibrisURL(searchString)

	response := utils.GetHTML(adlibrisURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	ScrapeAdlibrisPageData(document, searchString)
	fmt.Println("ADLIBRIS RESPONSE COMPLETE!")
}
