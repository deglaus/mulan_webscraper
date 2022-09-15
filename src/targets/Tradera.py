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


def traderaScraper(keyword):
    """ Scrapes items on Tradera.com, where keyword is the search term.

    :param keyword: the search term used to search Tradera
    :type keyword: string
    """
    # Load the HTML page
    driver.get(generateTraderaURL(keyword))
    table_name = "Tradera_" + keyword.replace(" ", "_")

    # Parse processed webpage with BeautifulSoup
    parser = BeautifulSoup(driver.page_source, features="lxml")
    # Only get the first section of items
    itemRow = []
    findRow = parser.find("section", {"class": "row mb-4"})
    if (findRow):
        itemRow = findRow.prettify()
    else:
        db_insert(table_name, [])
        return

    # Parse the item-section for items.
    parser = BeautifulSoup(itemRow, features="lxml")

    # -------------------------Title and url---------------------------------
    itemCardTitles = parser.find_all("p", {"class": "item-card-title"})

    # Parse each found item and make sure it actually contains a title

    for card in itemCardTitles:
        parserTitleUrl = BeautifulSoup(card.prettify(), features="lxml")
        found = parserTitleUrl.find("a", {"class": "font-weight-normal"})
        if (found):
            foundTitle = found.get_text().lstrip().rstrip()
            foundURL = "https://www.tradera.com" + found['href']
            time.sleep(0.01)
            itemData = scrapeTraderaItemPage(foundURL)

            foundItemData = [(keyword, "Tradera", foundURL, itemData[2],
                              foundTitle, itemData[3], itemData[4], itemData[0])]
            db_insert(table_name, foundItemData)


def convertPrice(card):
    """ Takes an html-element containing a text representing of a number and returns it as a float.

    :param card: the html-element containing the text
    :type: BeautifulSoup-object
    :rtype: float
    :return: The price within card converted to a float.    
    """
    cardText = card.get_text().lstrip().rstrip().replace(" ", " ")
    num = float(''.join(filter(str.isdigit, cardText)))
    return num


def scrapeTraderaItemPage(url):
    """ Scrapes data from a Tradera-item page and returns it in a list

    :param url: the url to the Tradera-item page
    :type url: string
    :rtype: list
    :return: A list containing a Tradera item's category, username of seller, url to item-picture and, price description. All of which as strings except for the price which is a float.
    """
    # Load the HTML page
    driver.get(url)

    # Parse processed webpage with BeautifulSoup
    parser = BeautifulSoup(driver.page_source, features="lxml")
    # Only get the first section of items
    category = parser.find(
        "span", {"class": "text-inter-medium ml-1 text-primary"}).get_text()

    seller = parser.find(
        "p", {"class": "size-london mb-0 text-truncate text-styled font-weight-bold seller-alias"}).get_text()

    picture = parser.find("img", {"class": "hover-img"})
    pictureURL = picture['src']

    description = parser.find(
        "div", {"class": "position-relative description mb-md-4"}).get_text()

    priceElement = parser.find(
        "div", {"class": "bid-details-amount"})
    if (not priceElement):
        priceElement = parser.find(
            "p", {"class": "bid-details-amount"})
    price = convertPrice(priceElement)

    return [category, seller, pictureURL, description, price]


def storeItemPageData(itemData, collectedItemData):
    """ Returns a list, collectedItemData with the contents in itemData added.

    :param itemData: the data to be added
    :type itemData: list
    :param collectedItemData: the list that data will be added to
    :type collectedItemData: list
    :rtype: list
    :return: Returns collectedItemData but with each element in itemData appended to the corresponding ordered element in collectedItemData.
    """
    amountOfData = len(itemData)

    for i in range(amountOfData):
        collectedItemData[i].append(itemData[i])

    return collectedItemData


def generateTraderaURL(searchString):
    """Returns a url to a Tradera search page for a search term

    :param searchString: the search term to search Tradera with
    :type searchString: string
    :rtype: string
    :return: A valid url to a search page on Tradera where the argument was used as the search term.
    """
    concatenatedString = "https://www.tradera.com/search?q=" + searchString
    return concatenatedString.replace(" ", "%20")


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
                     " (Id int unsigned NOT NULL AUTO_INCREMENT, SearchString VARCHAR(255), Site VARCHAR(255), URL VARCHAR(255), pictureURL VARCHAR(255), Title VARCHAR(255), Description text, Price float, Category VARCHAR(255), PRIMARY KEY(Id))")
    mycursor.execute("SHOW TABLES")

    sql = "INSERT INTO "+table_name + \
        " (SearchString, Site, URL, PictureURL, Title, Description, Price, Category) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)"
    vals = values

    mycursor.executemany(sql, vals)
    mydb.commit()


def scrapeAndInsertTradera(searchString):
    """ Scrapes a Tradera search result from a search term and stores into a local MySQL database.

    :param searchString: the search term to search Tradera with
    :type searchString: string
    """
    traderaScraper(searchString)
    print("TRADERA RESPONSE COMPLETE!")
