# Install google chrome, chromedriver and python3 (with pip).
#
# Pip-installations:
# pip install selenium
# pip install beautifulsoup4
# pip install mysql-connector-python
# pip install lxml
#
# If it tells you that webdriver manager is needed run:
# pip install webdriver-manager
#
# Chromedriver installation (match with your current version of Google Chrome):
# https://chromedriver.chromium.org/downloads
#
#
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager

from bs4 import BeautifulSoup
import os  # used for getting current directory, currently unused
import time

# Source: https://www.w3schools.com/python/python_mysql_insert.asp
import mysql.connector

# Instantiate options
opts = Options()
opts.add_argument("--headless")

# Instantiate a webdriver, install the latest version automatically.
driver = webdriver.Chrome(options=opts, service=Service(
    ChromeDriverManager().install()))

def scrollPortion(denominator):
    """ Execute a scroll script on a page in which the webdriver scrolls up and down.

    :param denominator: the fraction of which the initial scroll will be by. If denominator = 5 then 1/5 of the page will be scrolled to initially.
    :type denominator: int
    """
    current = driver.execute_script("return window.scrollY")
    driver.execute_script('window.scrollTo('+ str(current)  +', document.body.scrollHeight/' + str(denominator) + ')')
    time.sleep(0.3)
    current = driver.execute_script("return window.scrollY")
    distanceFromBottom = 999999
    driver.execute_script('window.scrollTo(document.body.scrollHeight/' + str(denominator)  + ', document.body.scrollHeight-'+ str(distanceFromBottom)  +')')
    driver.execute_script('window.scrollTo(document.body.scrollHeight-'+ str(distanceFromBottom)  +',' + str(current) + ' )')

def multipleScrolls(i):
    """ Will scroll the page up and down i/2 rounded to an integer times.

    :param i: double the amount of times the page will be scrolled.
    :type i: int
    """
    repetetor = int(i/2)
    for e in range(repetetor):
        while(i > 0):
            scrollPortion(i)
            i -= 1


def blocketScraper(keyword):
    """ Scrapes items on Blocket.se, where keyword is the search term.

    :param keyword: the search term used to search Blocket
    :type keyword: string
    """
    # Load the HTML page
    driver.get(generateBlocketURL(keyword))
    table_name = "Blocket_" + keyword.replace(" ", "_")
    
 #   button = driver.find_element_by_css_selector('#accept-ufti')
 #   button.click()
    # Scroll down page until it finishes loading images
    pageReady = False
    while(pageReady == False):
        button = driver.find_element_by_css_selector('#accept-ufti')
        button.click()
                
        pageReady = driver.execute_script('return document.readyState') == 'complete'


    multipleScrolls(10)

    # Parse processed webpage with BeautifulSoup
    parser = BeautifulSoup(driver.page_source, features="lxml")
    # Only get the first section of items
    itemRow = []
    findRow = parser.find("div", {
        "class": "SearchResults__SearchResultsWrapper-sc-10sdzls-0 bsVhZm"})
    if (findRow):
        itemRow = findRow.prettify()
    else:
        db_insert(table_name, [])
        return

    # Parse the item-section for items.
    parser = BeautifulSoup(itemRow, features="lxml")

    # -----------------------------Prices------------------------------------
    # Parse for prices
    prices = []
    itemCardPrices = parser.find_all(
        "div", {"class": "Price__StyledPrice-sc-1v2maoc-1 dNwEA"})

    
    for card in itemCardPrices:
        if (card.get_text() and card.get_text() != ""):
            prices.append(convertPrice(card))
        elif (card.get_text() == ""):
            prices.append(0)
        else:
            prices.append(0)

    # -------------------------Pictures------------------------------------
    pictureurls = []
    for pic in parser.select('img[class*="ListImage__"]'):
        foundPicUrl = pic['src']
        pictureurls.append(foundPicUrl)

    # -------------------------Remaining data---------------------------------
    parser = BeautifulSoup(itemRow, features="lxml")

    itemCardTitles = parser.find_all(
        "a", {"class": "Link-sc-6wulv7-0 styled__StyledTitleLink-sc-1kpvi4z-11 kVcpUt kbgaQK"})

    # Parse each found item to make sure it actually contains a title
    i = 0
    for card in itemCardTitles:
        if (card):
            foundTitle = card.get_text().lstrip().rstrip()
            foundURL = "https://www.blocket.se" + card['href']
            time.sleep(0.05)
            itemData = scrapeBlocketItemPage(foundURL)
            pictureurls = pictureAppender(pictureurls, i)

            foundItemData = [(keyword, "Blocket", foundURL, pictureurls[i],
                              foundTitle, itemData[2], prices[i], itemData[0])]
            db_insert(table_name, foundItemData)
            i += 1


def pictureAppender(pictureurls, i):
    """ Appends an empty string to a list if its length is less than a number i+3.
    :param pictureurls: the list
    :type: list
    :param i: the number
    :type: int
    :rtype: list
    :return: pictureurls with an empty string appended to its end if i+3 is greater than its original length.

    """

    if (len(pictureurls) < i+3):
        pictureurls.append("")
    return pictureurls


def convertPrice(card):
    """ Takes an html-element containing a text representing of a number and returns it as a float.

    :param card: the html-element containing the text
    :type: BeautifulSoup-object
    :rtype: float
    :return The price within card converted to a float, the float will be 0 if a price does not exist.    
    """

    cardText = card.get_text().lstrip().rstrip().replace(" ", "")
    if (cardText == ""):
        return 0.0

    num = float(''.join(filter(str.isdigit, cardText)))
    return num


def scrapeBlocketItemPage(url):
    """ Scrapes data from a Blocket-item page and returns it in a list

    :param url: the url to the Blocket-item page
    :type url: string
    :rtype: list
    :return: A list containing a Blocket item's category, username of seller, url to item-picture and description. All of which as strings.
    """

    driver.get(url)
    time.sleep(0.4)
    parser = BeautifulSoup(driver.page_source, features="lxml")

    categories = parser.find_all("div", {
        "class": "TextCallout2__TextCallout2Wrapper-sc-1bir8f0-0 dVkfPB Breadcrumb__ClickableBreadcrumb-sc-1ygccid-1 jryHxz"})
    if (categories):
        category = get_category(categories)
    else:
        category = "Alla"

    seller = parser.find("div", {
        "class": "TextSubHeading__TextSubHeadingWrapper-sc-1c6hp2-0 gQCEZy styled__AdvertiserName-sc-1f8y0be-7 epkISU"})
    foundSeller = trimSellerName(ifExists(seller))

    description = parser.find("div", {
        "class": "TextBody__TextBodyWrapper-cuv1ht-0 lieIzz BodyCard__DescriptionPart-sc-15r463q-2 gLpiBo"})
    foundDesc = ifExists(description)

    return [category, foundSeller, foundDesc]


def ifExists(obj):
    """ Returns a beautifulsoup-object's text if it is not a NoneType. Otherwise, an empty string is returned.

    :param obj: the beautiful soup object
    :type: object
    :rtype: string
    :return: obj's text if it is not a NoneType, otherwise an empty string.
    """

    if (obj):
        return obj.get_text()
    else:
        print("Didn't find obj")
        print(obj)
        return ""


def trimSellerName(longName):
    """ Removes unnecessary information from a Blocket text containing the seller's name.

    :param longName: the text from Blocket
    :type: string
    :rtype: string
    :return: longName except without anything after a '+' or the sub-strings 'Verifierad' or 'På Blocket'.
    """

    name = ""
    processedName = longName.replace(
        "Verifierad", "+").replace("På Blocket", "+")
    for char in processedName:
        if (char == "+"):
            return name
        else:
            name += char
    return name


def get_category(categories):
    """  Returns the second last element's text in a list of BeautifulSoup objects

    :param categories: contains BeautifulSoup objects containing category-names from Blocket
    :type: list
    :rtype: string
    :return: the second text of the second last element in categories.
    """
    deepestCategoryIndex = len(categories)-1
    return categories[deepestCategoryIndex].get_text()


def generateBlocketURL(searchString):
    """Returns a url to a Blocket search page for a search term

    :param searchString: the search term to search Blocket with
    :type searchString: string
    :rtype: string
    :return: A valid url to a search page on Blocket where the argument was used as the search term.
    """

    concatenatedString = "https://www.blocket.se/annonser/hela_sverige?q=" + searchString
    return concatenatedString.replace(" ", "+")


def db_insert(table_name, values):
    """ Inserts values into a table in a MySQL database.

    :param table_name: the name of the table
    :type tablename: string
    :param values: the values to be inserted into the table
    :type values: 
    """
    mydb = mysql.connector.connect(
        host="localhost",
        user="admin",
        password="mulan",
        database="test"
    )

    mycursor = mydb.cursor(buffered=True)
    mycursor.execute("CREATE TABLE IF NOT EXISTS " + table_name +
                     " (Id int unsigned NOT NULL AUTO_INCREMENT, SearchString VARCHAR(255), Site VARCHAR(255), URL VARCHAR(255), PictureURL VARCHAR(255), Title VARCHAR(255), Description text, Price float, Category VARCHAR(255), PRIMARY KEY(Id))")
    mycursor.execute("SHOW TABLES")

    sql = "INSERT INTO "+table_name + \
        " (SearchString, Site, URL, PictureURL, Title, Description, Price, Category) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)"
    vals = values

    mycursor.executemany(sql, vals)
    mydb.commit()


def scrapeAndInsertBlocket(searchString):
    """ Scrapes a Blocket search result from a search term and stores into a local MySQL database.

    :param searchString: the search term to search Blocket with
    :type searchString: string
    """

    blocketScraper(searchString)
    print("BLOCKET RESPONSE COMPLETE!")
