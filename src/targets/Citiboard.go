package targets

import (
	"fmt"
	"time"
	// "regexp"
	"github.com/go-rod/rod"
	"gosecondhand/src/database"
	shsutils "gosecondhand/src/utils"
	"strconv"
	"strings"
	//"github.com/go-rod/rod/lib/cdp"
	//"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	//"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"regexp"
	"sync"
)

// TODO:
// Citiboard has a very elaborate and accurate location system
// You can specify even the individual city in a search
// This will no doubt be relevant to users, but will also require some menial work
func GenerateCitiboardURL(searchString string) string {
	searchString = strings.ReplaceAll(searchString, " ", "+")
	// For now, we make a search for the entire country
	concatenatedString := ("https://citiboard.se/hela-sverige?inp_text=" + searchString)
	return concatenatedString
}

// Figure out max pages
// Figure out max offset pased on max pages
// Check if results are empty

// Figure out number of items on a page
// Scrape those items
// Go to next page
// Figure out number of items on a page
// Scrape those items
// Go to next page
// ...

func CitiboardResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()
	///////////////////
	// DEBUG OPTIONS //
	///////////////////
	l := launcher.New().
		Headless(true).
		Devtools(true)
	defer l.Cleanup()
	url := l.MustLaunch()
	browser := rod.New().
		ControlURL(url).
		Trace(false).
		//SlowMotion(500 * time.Millisecond). // set delay
		MustConnect()
	launcher.Open(browser.ServeMonitor(""))
	//////////////
	// SCARPING //
	//////////////
	defer browser.MustClose()
	// Here we generate our own URL (unlike e.g, Etsy) because
	// the user will perhaps want to specify a city
	CitiboardURL := GenerateCitiboardURL(searchString)
	page := browser.MustPage(CitiboardURL).MustWindowMaximize()

	db := database.ConnectDB()
	defer db.Close()
	err := db.Ping()
	shsutils.CheckIfError(err)
	// Create tables
	table_name := "Citiboard_" + strings.ReplaceAll(searchString, " ", "_")

	// Get total number of results
	// Could be incorporated in a major refactoring, lazy fix
	totalResultsText := page.MustElement(".searchResultsHeader").MustText()
	reTotalResults := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	reTotalResultsResult := reTotalResults.FindAllString(totalResultsText, -1)
	totalResultsNumbersOnly := reTotalResultsResult[0]
	totalResultsNumbersOnly = strings.ReplaceAll(totalResultsNumbersOnly, ",", "")
	totalResults, _ := strconv.Atoi(totalResultsNumbersOnly)
	

	currentArticle := 1
	offsetMod := 0
	totalScraped := 0
	currentURL := CitiboardURL + "&offset=" + strconv.Itoa(offsetMod)
	for i := 1; i <= totalResults; i++ {
		// DESCRIPTION
		description := page.MustElement(".boardItems > .masonry-layout > article:nth-child(" + strconv.Itoa(currentArticle+1) + ") > div:nth-child(2) > .gridFooter > .gridLocation").MustText()
		// URL
		URLRaw := page.MustElement(".boardItems > .masonry-layout > article:nth-child(" + strconv.Itoa(currentArticle+1) + ") > div:nth-child(2) > .gridTitle > h3 > a").MustProperty("href")
		URLsuffix := URLRaw.JSON("> ", " ")
		URL := strings.Trim(URLsuffix, "\"")
		// PICTURE URL
		picURLSelector := ".boardItems > .masonry-layout > article:nth-child(" + strconv.Itoa(currentArticle+1) + ") > div:nth-child(2) > .picture > figure > a > img"
		var picURL string
		if page.MustHas(picURLSelector) {
			picURLRaw := page.MustElement(picURLSelector).MustProperty("src")
			picURL = picURLRaw.JSON("> ", " ")
		} else {
			picURL = ""
		}
		// TITLE
		title := page.MustElement(".boardItems > .masonry-layout > article:nth-child(" + strconv.Itoa(currentArticle+1) + ") > div:nth-child(2) > .gridTitle > h3").MustText()
		// PRICE
		var price float64
		priceSelector := ".boardItems > .masonry-layout > article:nth-child(" + strconv.Itoa(currentArticle+1) + ") > div:nth-child(2) > .gridPrice"
		if page.MustHas(priceSelector) {
			priceText := page.MustElement(priceSelector).MustText()
			rePrice := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			rePriceResult := rePrice.FindAllString(priceText, -1)
			priceNumbersOnly := rePriceResult[0]
			priceNumbersOnly = strings.ReplaceAll(priceNumbersOnly, ",", "")
			price, _ = strconv.ParseFloat(priceNumbersOnly, 64)

		} else {
			price = 0
		}

		// Insert
		data := shsutils.Item{}
		data.SearchString = searchString
		data.Site = "Citiboard"
		data.Description = description
		data.URL = URL
		data.PictureURL = picURL
		data.Title = title
		data.Price = price
		database.InsertData(db, table_name, data)

		currentArticle++
		totalScraped++

		if totalScraped > 50 {
			fmt.Println("CITIBOARD RESPONSE COMPLETE!")
			return
		}

		if currentArticle > 20 {
			offsetMod = offsetMod + 23
			currentArticle = 1

			currentURL = CitiboardURL + "&offset=" + strconv.Itoa(offsetMod)
			page.MustNavigate(currentURL)
			time.Sleep(2 * time.Second)
		}
	}

	fmt.Println("CITIBOARD RESPONSE COMPLETE!")
	return
	utils.Pause()
}
