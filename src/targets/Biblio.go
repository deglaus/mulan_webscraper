// BUG: Something is wrong with the float price ocassionally (gives back like ~20 decimal points)

// Put this code in SQL to make sure to avoid "incorrect string value":
// ALTER DATABASE test CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
package targets

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosecondhand/src/database"
	"gosecondhand/src/utils"
	//"log"
	"strconv"
	"strings"
	"sync"
	//"math"
	"time"
	"unicode/utf8"
)

func getTotalPageNumBiblio(biblioSearchId string) int {
	response := utils.GetHTML(biblioSearchId)
	defer response.Body.Close()
	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)
	// Default total page num to 1
	foundTotalPageNum := "1"
	foundTotalPageNumInt := 1
	// Find actual total page num
	foundTotalPageNum = document.Find("ul.pagination").Text()
	if len(foundTotalPageNum) > 1 {
		temp := strings.Split(foundTotalPageNum, "\n")
		// Assume first is max
		maxValue, _ := strconv.Atoi(temp[0])
		// Convert each value from string to int and check for max
		for _, value := range temp {
			if valueInt, err := strconv.Atoi(value); err == nil && valueInt > maxValue {
				maxValue = valueInt
			}
		}
		foundTotalPageNumInt = maxValue
	}
	//fmt.Println("Total number of pages:", foundTotalPageNumInt)
	if foundTotalPageNumInt > 10 {
		foundTotalPageNumInt = 10
	}
	return foundTotalPageNumInt
}

func GenerateBiblioSearchId(searchString string) string {
	searchString = strings.ReplaceAll(searchString, " ", "+")
	resultsPerPageListing := 50
	concatenatedString := ("https://www.biblio.com/search.php?stage=1&author=&title=&isbn=&keyisbn=" + searchString + "&publisher=&illustrator=&mindate=&maxdate=&minprice=&maxprice=&country=&format=&cond=&days_back=&order=priceasc&pageper=" + strconv.Itoa(resultsPerPageListing) + "&quantity=")

	response := utils.GetHTML(concatenatedString)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	
	var foundSearchIdURLFixed string
	foundSearchId := document.Find("li.next>a")
	foundSearchIdURL, _ := foundSearchId.Attr("href")
	if (len(foundSearchIdURL) > 0) {
		foundSearchIdURLFixed = strings.ReplaceAll(foundSearchIdURL, "page=2", "page=1") // Lazy, fix.		
	} else {
		foundSearchIdURLFixed = concatenatedString
	}
	

	fmt.Println("URL with found search ID:", foundSearchIdURLFixed)

	return foundSearchIdURLFixed
}

func GenerateBiblioURL(searchIdURL string, pageNum int) string {
	searchIdURL = strings.ReplaceAll(searchIdURL, "page=1", "page="+strconv.Itoa(pageNum))
	return searchIdURL
}

func acquireAndInsertBiblioData(db *sql.DB, table_name string, item *goquery.Selection, searchString string) {
	data := utils.Item{}
	// URL //
	foundURL := item.Find("h2.title>a")
	if len(foundURL.Text()) < 1 {
		return
	}
	URL, _ := foundURL.Attr("href")
	//fmt.Println(URL)
	// PICTURE URL //
	foundPicURL := item.Find("div.image>img")
	picURL, _ := foundPicURL.Attr("data-pagespeed-lazy-src")
	//fmt.Println(picURL)
	// TITLE //
	foundTitle := item.Find("h2.title>a")
	title := foundTitle.Text()
	//fmt.Println(title)
	// DESCRIPTION //
	descriptionRawBasic := item.Find("dl.fact-list").Text()
	descriptionRawFull := item.Find("div.item-description").Text()
	description := descriptionRawBasic + descriptionRawFull
	description = strings.Join(strings.Fields(description), " ")
	//fmt.Println("description:", description)
	// PRICE //
	// Base
	priceBaseFound := item.Find("span.item-price").Text()
	priceBaseParts := strings.Split(priceBaseFound, "US$")
	if len(priceBaseParts) > 1 {
		priceBaseParts = strings.Split(priceBaseParts[1], " ")
	}
	priceBase, err := strconv.ParseFloat(priceBaseParts[0], 32)
	utils.CheckIfError(err)
	//fmt.Println("PriceBase:", priceBase)
	// Shipping
	priceShipFound := item.Find("div.shipping").Text()
	priceShipParts := strings.Split(priceShipFound, "US$")
	if len(priceShipParts) > 1 {
		priceShipParts = strings.Split(priceShipParts[1], "shipping to SWE")
		priceShipParts[0] = strings.TrimSuffix(priceShipParts[0], "\n")
	}
	priceShip, err := strconv.ParseFloat(priceShipParts[0], 32)
	utils.CheckIfError(err)
	//fmt.Println("PriceShipping:", priceShip)
	// Total
	totalPrice := priceBase + priceShip
	// DB INSERTION //
	data.SearchString = searchString
	data.Site = "Biblio"

	// https://henvic.dev/posts/go-utf8/
	// We need to do something about standardizing encodings
	// for now, we just ignore encodings that are incompatible with db
	//fmt.Println(utf8.ValidString(description))
	if utf8.ValidString(description) {
		data.Description = description
	}

	data.URL = URL
	data.PictureURL = picURL

	//fmt.Println(utf8.ValidString(title))
	if utf8.ValidString(title) {
		data.Title = title
	}

	data.Price = totalPrice
	database.InsertData(db, table_name, data)

}

func singlePageResponseBiblio(biblioURL string, searchString string) {
	response := utils.GetHTML(biblioURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	// Connect to database
	db := database.ConnectDB()
	defer db.Close()
	err := db.Ping()
	utils.CheckIfError(err)
	table_name := "Biblio_" + strings.ReplaceAll(searchString, " ", "_")

	document.Find("div.item").Each(func(index int, item *goquery.Selection) {
		acquireAndInsertBiblioData(db, table_name, item, searchString)
	})
}

func BiblioResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()

	biblioSearchId := GenerateBiblioSearchId(searchString)
	//fmt.Println(biblioSearchId)
	totalPageNum := getTotalPageNumBiblio(biblioSearchId)
	//fmt.Println(totalPageNum)

	for i := 1; i <= totalPageNum; i++ {
		biblioURL := GenerateBiblioURL(biblioSearchId, i)
		fmt.Println(biblioURL)
		singlePageResponseBiblio(biblioURL, searchString)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("BIBLIO RESPONSE DONE!")
}
