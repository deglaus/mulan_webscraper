package targets

import (
	"fmt"
	//"time"
	"github.com/go-rod/rod"
	"gosecondhand/src/database"
	shsutils "gosecondhand/src/utils"
	"regexp"
	"strconv"
	"strings"
	//"github.com/go-rod/rod/lib/cdp"
	//"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	//"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"sync"
	"time"
)

func pickCity(city string) string {
	var cityUrlString string
	switch city {
	case "gothenburg":
		cityUrlString = "gothenburg"
	case "malmö":
		cityUrlString = "110837332277562"
	case "uppsala":
		cityUrlString = "110976692260411"
	case "upplands väsby":
		cityUrlString = "106931919342119"
	case "västerås":
		cityUrlString = "110873668941216"
	case "örebro":
		cityUrlString = "110611878960213"
	case "linköping":
		cityUrlString = "110680185620981"
	case "helsingborg":
		cityUrlString = "110721515622776"
	case "jönköping":
		cityUrlString = "110604105633625"
	case "norrköping":
		cityUrlString = "105830426123838"
	case "lund":
		cityUrlString = "108679122496232"
	case "umeå":
		cityUrlString = "110664898961620"
	case "gävle":
		cityUrlString = "114779245206224"
	default:
		cityUrlString = "stockholm"
	}

	return cityUrlString
}

func GenerateFacebookURL(searchString string, city string) string {
	cityURL := pickCity(city)
	concatenatedString := (`https://www.facebook.com/marketplace/` + cityURL + `/search/?query=` + searchString)
	return concatenatedString
}

func FacebookMarketResponse(searchString string, city string, wg *sync.WaitGroup) {
	defer wg.Done()

	/////////////////////
	// BROWSER OPTIONS //
	/////////////////////
	l := launcher.New().
		Headless(true).
		Devtools(false) // This causes some sites to change structure
	defer l.Cleanup()
	url := l.MustLaunch()
	browser := rod.New().
		ControlURL(url).
		Trace(false).
		//SlowMotion(10 * time.Millisecond). // set delay
		MustConnect()
	launcher.Open(browser.ServeMonitor(""))
	////////////////////
	// ERROR HANDLING //
	////////////////////
	// check := func(err error) {
	// 	var evalErr *rod.ErrEval
	// 	if errors.Is(err, context.DeadlineExceeded) { // timeout error
	// 		fmt.Println("timeout err")
	// 	} else if errors.As(err, &evalErr) { // eval error
	// 		fmt.Println(evalErr.LineNumber)
	// 	} else if err != nil {
	// 		fmt.Println("can't handle", err)
	// 	}
	// }
	//////////////
	// SCRAPING //
	//////////////
	defer browser.MustClose()
	// Here we generate our own URL (unlike e.g, Etsy) because
	// the user will perhaps want to specify a city
	page := browser.MustPage(GenerateFacebookURL(searchString, city)).MustWindowMaximize()
	// Press accept in cookies dialogue window
	page.Timeout(30*time.Second).MustElementR("span", "Only allow essential cookies").MustClick()

	// CHECK: NO RESULTS? //
	if page.MustHasR("span", "Try a new search. Check spelling, change your filters, or try a less specific search term.") {
		fmt.Println("No search results for Facebook!")
		return
	}

	time.Sleep(5 * time.Second) // Ensure site has fully loaded
	// EXPLANATION:
	// class div.f9o22wc5:nth-child(<even number>) (excluding 0)
	// is a separator and should be ignored
	//
	// class div.f9o22wc5:nth-child(<uneven num>)
	// contains relevant content and shouldn't be ignored
	totalScraped := 0
	itemCounter := 1
	itemChunkCounter := 1
	for {
		if itemCounter > 12 {
			itemCounter = 1
			itemChunkCounter = itemChunkCounter + 2
		}
		// URL
		URLRaw := page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ") > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1)").MustScrollIntoView().MustProperty("href")
		// For loading new listings
		//time.Sleep(200 * time.Millisecond)
		URL := URLRaw.JSON("> ", " ")
		URL = strings.Trim(URL, "\"")

		// TITLE
		title := page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ") > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(2)").MustText()

		// PRICE
		priceText := page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ") > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > span:nth-child(1)").MustText()
		var price float64
		rePrice := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		rePriceResult := rePrice.FindAllString(priceText, -1)
		if len(rePriceResult) != 0 {
			priceNumbersOnly := rePriceResult[0]
			priceNumbersOnly = strings.ReplaceAll(priceNumbersOnly, ",", "")
			price, _ = strconv.ParseFloat(priceNumbersOnly, 64)
		} else {
			price = 0
		}

		// PICTURE URL
		picURLRaw := page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ") > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > img:nth-child(1)").MustProperty("src")
		picURL := picURLRaw.JSON("> ", " ")
		picURL = strings.Trim(picURL, "\"")

		// DESCRIPTION (location)
		description := page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ") > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(3) > span:nth-child(1) > div:nth-child(1) > span:nth-child(1)").MustText()

		//page.MustNavigate(URL)
		// page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ")").MustClick()
		////////////////////////////
		// INSIDE INDIVIDUAL PAGE //
		////////////////////////////
		// time.Sleep(750 * time.Millisecond)
		// for (page.MustHasR("span", "Please log in to see this page.")) {
		// 	page.MustNavigate(URL)
		//  	time.Sleep(5 * time.Second)
		// page.MustElement(".cypi58rs").MustClick()
		// page.MustElement(".gh1tjcio").MustClick()
		// page.MustElement("div.byvelhso:nth-child(1) > div:nth-child(1)").MustClick()
		// page.Timeout(30 * time.Second).MustElement("div.f9o22wc5:nth-child(" + strconv.Itoa(itemChunkCounter) + ") > div:nth-child(2) > div:nth-child(" + strconv.Itoa(itemCounter) + ")").MustClick()

		//}
		// <NOT YET USED> RETRIEVAL OF ALL PICTURE URLS
		// var allPicLinks string
		// allPics := page.Timeout(30 * time.Second).MustElements("img.datstx6m")
		// for i, pic := range allPics {
		// 	picLinkRaw := pic.MustProperty("src")
		// 	//fmt.Println("index:", i, picLink)
		// 	picLink := picLinkRaw.JSON("> ", " ")
		// 	picLink = strings.Trim(picLink, "\"")
		// 	allPicLinks = allPicLinks + " " + strconv.Itoa(i) + " " + picLink
		// 	fmt.Println("ALL PICTURES: ", allPicLinks)
		// }
		// // PICTURE RUL
		// picURLRaw := page.Timeout(30 * time.Second).MustElement("img.datstx6m").MustProperty("src")
		// picURL := picURLRaw.JSON("> ", " ")
		// picURL = strings.Trim(picURL, "\"")

		// DESCRIPTION
		// Regular description
		// var description string
		// if (page.MustHasR("span", "Seller's Description")) {
		// 	description = page.Timeout(30 * time.Second).MustElement(".ii04i59q > div:nth-child(1) > span:nth-child(1)").MustText()

		// } else if (page.MustHasR("span", "Description")) {
		// 	description = page.Timeout(30 * time.Second).MustElement("div.n851cfcs:nth-child(8) > div:nth-child(2)").MustText()
		// } else {
		// 	description = page.Timeout(30 * time.Second).MustElement("div.muag1w35:nth-child(1) > div:nth-child(2)").MustText()
		// }

		// Screenshot of location
		// NOTE: This is not meant for storage in the database
		// Only directly downloaded to some filepath in CDN
		// and then delivered to user
		// page.Timeout(30 * time.Second).MustElement("div.oajrlxb2").MustScreenshot("my.png")

		// DB INSERTION
		// Connect to database
		db := database.ConnectDB()
		defer db.Close()
		err := db.Ping()
		shsutils.CheckIfError(err)
		// Create tables
		table_name := "FacebookMarket_" + strings.ReplaceAll(searchString, " ", "_")
		// Insert
		data := shsutils.Item{}
		data.SearchString = searchString
		data.Site = "FacebookMarket"
		data.Description = description
		data.URL = URL
		data.PictureURL = picURL
		data.Title = title
		data.Price = price
		database.InsertData(db, table_name, data)

		//page.NavigateBack()
		//time.Sleep(5 * time.Second)

		itemCounter++
		totalScraped++
		if totalScraped > 60 {
			fmt.Println("FB RESPONSE COMPLETE!")
			return
		}
	}
	utils.Pause()
}

// div.f9o22wc5:nth-child(1) > div:nth-child(2)
// div.f9o22wc5:nth-child(1) > div:nth-child(2) > div:nth-child(1)
// div.f9o22wc5:nth-child(1) > div:nth-child(2) > div:nth-child(2)
// div.f9o22wc5:nth-child(1) > div:nth-child(2) > div:nth-child(3)
// div.f9o22wc5:nth-child(Y) > div:nth-child(2) > div:nth-child(X) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1)
// ...

// div.f9o22wc5:nth-child(3) > div:nth-child(2) > div:nth-child(1)
