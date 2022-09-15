package targets

import (
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	"gosecondhand/src/database"
	shsutils "gosecondhand/src/utils"
	//"log"
	"github.com/go-rod/rod"
	"math"
	"strconv"
	"strings"
	"time"
	//"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	//"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	//"github.com/ysmood/gson"
	"regexp"
	"sync"
	//"reflect"
)

func GenerateEtsyURL(searchString string, pageNum int) string {
	searchString = strings.ReplaceAll(searchString, " ", "+")
	concatenatedString := ("https://www.etsy.com/search?q=" + searchString + "&page=" + strconv.Itoa(pageNum))
	fmt.Println("Generated search string:", concatenatedString)
	return concatenatedString
}

func EtsyResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done()
	///////////////////
	// DEBUG OPTIONS //
	///////////////////
	// Headless runs the browser on foreground, you can also use env "rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		Headless(true).
		Devtools(false)
	defer l.Cleanup()
	url := l.MustLaunch()
	// Trace shows verbose debug information for each action executed
	// Slowmotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(false).
		//SlowMotion(100 * time.Millisecond). // set delay
		MustConnect()
	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with env rod=monitor
	launcher.Open(browser.ServeMonitor(""))
	//////////////
	// SCARPING //
	//////////////
	// Initialize variables (needed since now inside if-statement)
	var stars string
	var seller string
	var URL string
	var picURL string
	var title string
	var price float64
	var description string

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()
	// Create a new page
	page := browser.MustPage("https://etsy.com").MustWindowMaximize()
	// Find accept cookie button and click it
	page.Timeout(30 * time.Second).MustElement(".wt-btn.wt-btn--filled.wt-mb-xs-0").MustClick()
	// INPUT searchString in "input" and press enter
	page.Timeout(30 * time.Second).MustElement("input").MustInput(searchString).MustPress(input.Enter)
	// CHECK: NO RESULTS? //
	// Be aware, might cause bugs where this div
	// as described by css selector is used elsewhere
	if page.MustHas("div.wt-position-relative:nth-child(4)") {
		fmt.Println("No search results for Etsy!")
		return
	}
	// TOTAL RESULTS //
	resultsNumRawText := page.Timeout(30 * time.Second).MustElement("div.wt-text-right-xs").MustText()

	// Check if it's a different vairant of the page. And if it is,
	// Update a sequence modifier for iteration due to 4 extra divs at top
	// This is only for the first page however, so it will be reset at
	// the bottom of the loop
	offsetMod := 0
	if strings.Contains(resultsNumRawText, "with Ads") {
		offsetMod = 4
	}

	reResults := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	reResultsResult := reResults.FindAllString(resultsNumRawText, -1)
	resultsNum := reResultsResult[0]
	resultsNum = strings.ReplaceAll(resultsNum, ",", "")
	resultsFloat, err := strconv.ParseFloat(resultsNum, 64)
	shsutils.CheckIfError(err)
	// NUMBER OF PAGES //
	// 48 results per page
	pageNum := resultsFloat / 48
	// Separate fraction from integer
	pageNumFullPages := math.Floor(pageNum)
	// RESULTS ON LAST PAGE //
	// Get fraction modifier
	lastPageItemAmountMod := pageNum - pageNumFullPages
	// Apply fraction modifier to 48 for results on last page
	lastPageItemAmount := math.Round(48 * lastPageItemAmountMod)
	// Type-cast rounded down float to integer
	pageNumFullPagesInt := int(pageNumFullPages)

	// CHECK: SURPASSES UPPER LIMIT? //
	// (2 is arbitrary)
	if pageNumFullPagesInt > 2 {
		pageNumFullPagesInt = 2
	}

	db := database.ConnectDB()
	defer db.Close()
	err = db.Ping()
	shsutils.CheckIfError(err)
	// Create tables
	table_name := "Etsy_" + strings.ReplaceAll(searchString, " ", "_")

	// PAGE-DETERMINING LOOP //
	for i := 1; i <= pageNumFullPagesInt+1; i++ {
		page.MustNavigate("https://www.etsy.com/se-en/search?q=" + searchString + "&page=" + strconv.Itoa(i) + "&ref=pagination")
		currentPageItemAmount := 48
		// If it's the last page, change to modified amount
		if i == pageNumFullPagesInt+1 {
			currentPageItemAmount = int(lastPageItemAmount)
		}
		// ON-PAGE WORK LOOP //
		for j := 1 + offsetMod; j <= currentPageItemAmount+offsetMod; j++ {
			hasItemCounter := 0
			// URL
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1)") {
				URLRaw := page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1)").MustProperty("href")

				URL = URLRaw.JSON("> ", " ")
				URL = strings.Trim(URL, "\"")

				hasItemCounter++
			}

			// PICTURE URL
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(1) > div > div > div > div > div > img") {
				picURLRaw := page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(1) > div > div > div > div > div > img").MustProperty("src")

				picURL = picURLRaw.JSON("> ", " ")
				picURL = strings.Trim(picURL, "\"")

				hasItemCounter++
			}

			// TITLE
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(2) > div > h3") {
				title = page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(2) > div > h3").MustText()

				hasItemCounter++
			}

			// DESCRIPTION
			// Stars
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(2) > div > div") {
				stars = page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div > div > a:nth-child(1) > div:nth-child(2) > div > div").MustText()

				hasItemCounter++
			}

			// Seller
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(5) > p:nth-child(1) > span:nth-child(3)") {
				seller = page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(5) > p:nth-child(1) > span:nth-child(3)").MustText()

				hasItemCounter++
			}

			// TODO: Discount
			// TODO: "only x left soon"
			// Combine description
			if hasItemCounter > 1 {
				description = stars + " ....... " + seller
			}
			//fmt.Println(description)
			// PRICE
			if page.MustHas(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(4) > p:nth-child(1)") {
				priceText := page.Timeout(30 * time.Second).MustElement(".tab-reorder-container > li:nth-child(" + strconv.Itoa(j) + ") > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(4) > p:nth-child(1)").MustText()
				rePrice := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
				rePriceResult := rePrice.FindAllString(priceText, -1)
				priceNumbersOnly := rePriceResult[0]
				priceNumbersOnly = strings.ReplaceAll(priceNumbersOnly, ",", "")
				price, err = strconv.ParseFloat(priceNumbersOnly, 64)
				shsutils.CheckIfError(err)

				hasItemCounter++
			}

			// DB INSERTION
			// Connect to database
			if hasItemCounter == 6 {

				// Insert
				data := shsutils.Item{}
				data.SearchString = searchString
				data.Site = "Etsy"
				data.Description = description
				data.URL = URL
				data.PictureURL = picURL
				data.Title = title
				data.Price = price
				database.InsertData(db, table_name, data)
			}

			offsetMod = 0
		}
	}

	fmt.Println("ETSY RESPONSE COMPLETE")
	return
	utils.Pause()
}
