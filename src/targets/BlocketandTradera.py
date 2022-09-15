import Tradera
import Blocket
from threading import Thread


def scrapeAndInsert(searchString):
    Thread(target=Blocket.scrapeAndInsertBlocket, args=[searchString]
           ).start()
    Thread(target=Tradera.scrapeAndInsertTradera, args=[searchString]).start()
