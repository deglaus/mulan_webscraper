///////////////////
// Prerequisites //
///////////////////
// --- HTTP server and MySQL required (e.g, Xampp) ---
// 1. Start MySQL and HTTP-server (Apache for XAMPP)
// 2. In your web browser, type "localhost/phpmyadmin" in the address bar
// 3. Click on the tab "user accounts", and then "add account"
// 4. Enter the "admin" as user name, "localhost" (without %) as hostname, and "mulan" as password.
// 5. Check all privilages, and click "Go"-button. You can now run the code.
// 6a. Go to "~/xampp/mysql/bin" or equivalent and type "mysql.exe --user=admin --password=mulan". You can now run MySQL commands.
// 6b. Alternatively, return to localhost/phpmyadmin to manipulate databse with gui.

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // underscore means we will import, but not directly reference
	"gosecondhand/src/utils"
	"log"
)

//////////////////////
// PUBLIC FUNCTIONS //
//////////////////////

// ConnectDB opens connection with a MySQL database and returns it.
// Note: This is not a connection in the networking sense, it is merely
// a "pathway" through which you can send queries.
//
// Example:
//         db := database.ConnectDB()
//         defer db.Close()
//         err := db.Ping()
//         if err != nil {
//                 ...
//         }
//
func ConnectDB() (db *sql.DB) {
	db_user := "admin"
	db_pass := "mulan"
	db_addr := "localhost" // host name
	db_dbas := "test"      // database name

	// Print out relevant information
	fmt.Println("Drivers:", sql.Drivers())
	s := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_user, db_pass, db_addr, db_dbas)
	fmt.Println(s)

	// Establish connection
	db, err := sql.Open("mysql", s)
	if err != nil {
		log.Fatal("Unable to open connection with DB: %w", err)
	}

	return db
}

// CreateTables takes a database connection and creates all needed tables for the scraper.
// Regardless of future changes to the system architecture, this is the function from which
// Universal table creation is performed.
//
// Example:
//         db := database.ConnectDB()
//         defer db.Close()
//         err := database.CreateTables(db)
//         if err != nil {
//                 ...
//         }
//
func CreateTables(db *sql.DB, table_name string) error {
	createTableGeneric(db, table_name)
	// createTableBlocket(db)
	// createTableEbay(db)
	// createTableFacebook(db)
	// createTableTradera(db)
	return nil
}

// DropAllTables takes a database, and drops all tables within its 'test' database
//
//
// Example:
//     db := database.ConnectDB()
//     ...
//     defer db.Close()
//     ...
//      err := db.Ping()
//      if err != nil {
//          log.Fatal("Connection could not be verified with Ping(): ", err)
//      }
//     ...
//	   err = database.DropAllTables(db)
//	   if err != nil {
//         log.Fatal("Something went wrong during table dropping call: ", err)
//	   }
func DropAllTables(db *sql.DB) error {
	names := getTableNames(db)

	for i := 0; i < len(names); i++ {

		table_name := names[i]
		query := "DROP TABLE IF EXISTS " + table_name + ";"

		drop, err := db.Prepare(query)
		defer drop.Close()
		if err != nil {
			log.Fatal("Dropping table: "+table_name+" went wrong in preparation stage: %w", err)
		}

		_, err = drop.Exec()
		if err != nil {
			log.Fatal("Dropping table: "+table_name+" went wrong in execution stage: %w", err)
		}
	}

	return nil
}

// IsInDatabase takes a database and a keyword and returns a boolean
// which indicates whether the keyword has been logged in the database
// in at least one table.
//
// Example:
//
//  db := database.ConnectDB()
//  defer db.Close()
//	err := db.Ping()
//	if err != nil {
//		log.Fatal("Connection could not be verified with Ping(): ", err)
//	}
//  ...
//  isAlreadyInDB := database.IsInDatabase(db, searchString)
func IsInDatabase(db *sql.DB, keyword string) bool {
	names := getTableNames(db)

	for i := 0; i < len(names); i++ {
		current_table_name := names[i]
		if parseTableName(current_table_name) == keyword {
			return true
		}
	}
	return false
}

// InsertData takes a database, the name of a table, an item, and inserts this item in
// given table in given databse. [SUBJECT TO CHANGE]
//
// Example:
//         db := database.connectDB()
//         ...
//         err := database.CreateTables(db)
//         ...
//         err = InsertData(db, "Generic", {...})
//         if err != nil {
//                 ...
//         }
//
func InsertData(db *sql.DB, tableName string, item utils.Item) error {

	query := "INSERT INTO " + tableName + " (SearchString, Site, URL, PictureURL, Title, Description, Price, Category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	insert, err := db.Prepare(query) // "Prepare" will check if query is a valid mysql-query
	defer insert.Close()
	if err != nil {
		log.Fatal("Insertion went wrong in preparation stage: ", err)
	}

	// Execute insertion
	// Note: the parameters to exec() correspond to question marks in query
	_, err = insert.Exec(item.SearchString, item.Site, item.URL, item.PictureURL, item.Title, item.Description, item.Price, item.Category)
	if err != nil {
		log.Fatal("Insertion went wrong in the execution stage: %w", err)
	}

	return nil
}

func GetAllData(db *sql.DB, tableName string) (items []utils.Item, err error) {

	// Note: .Query() instead of .Exec(), this is because we are now
	//       Asking for results, and not just executing changes
	response, err := db.Query("SELECT * FROM " + tableName)
	defer response.Close()
	if err != nil {
		return items, err
	}

	// Loop through response using .Next()
	for response.Next() {
		var tempItem utils.Item // Placeholder variable for item

		// Copies columns in current row into values pointed to in arguments
		// It essentially converts columns read from the database into common Go types
		// as well as special types provided by sql package
		err = response.Scan(&tempItem.Id, &tempItem.SearchString, &tempItem.Site, &tempItem.URL, &tempItem.PictureURL, &tempItem.Title, &tempItem.Description, &tempItem.Price, &tempItem.Category)
		if err != nil {
			return items, err
		}

		items = append(items, tempItem) // append temp item to array of items
	}

	return items, nil
}

func GetAllPriceData(db *sql.DB, tableName string) (itemPrices []float64, err error) {

	response, err := db.Query("SELECT Price FROM " + tableName)
	defer response.Close()
	if err != nil {
		return itemPrices, err
	}

	for response.Next() {
		var tempPrice float64
		err = response.Scan(&tempPrice)
		if err != nil {
			return itemPrices, err
		}
		itemPrices = append(itemPrices, tempPrice)
	}

	return itemPrices, nil
}

func ClearTable(db *sql.DB, tableName string) error {

	query := "TRUNCATE TABLE " + tableName + ""

	clear, err := db.Prepare(query)
	defer clear.Close()
	if err != nil {
		log.Fatal("Something went wrong in preperation stage of clear", err)
	}

	_, err = clear.Exec()
	if err != nil {
		log.Fatal("Something went wrong in execution stage of clear", err)
	}

	return nil
}

func DeleteAllAbovePrice(db *sql.DB, tableName string, priceLimit float64) error {
	query := "DELETE FROM " + tableName + " WHERE price > ?"

	del, err := db.Prepare(query)
	defer del.Close()
	if err != nil {
		return err
	}

	_, err = del.Exec(priceLimit)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePrice(db *sql.DB, tableName string, newPrice float64, itemID int) error {
	query := "UPDATE " + tableName + " SET price = ? WHERE Id like ?"

	update, err := db.Prepare(query)
	defer update.Close()
	if err != nil {
		return err
	}

	_, err = update.Exec(newPrice, itemID)
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////
// NON-PUBLIC FUNCTIONS //
//////////////////////////

// Returns all table names in the test-database
func getTableNames(db *sql.DB) []string {
	type Tag struct {
		Table_name string `db:"table_name" json:"table_name"`
	}
	names := []string{}
	names_query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'test';"

	results, err_names := db.Query(names_query)

	if err_names != nil {
		log.Fatal("Getting table names went wrong went wrong in preparation stage: %w", err_names)
	}

	for results.Next() {

		var tag Tag
		err_names = results.Scan(&tag.Table_name)
		if err_names != nil {
			panic(err_names.Error())
		}
		names = append(names, tag.Table_name)

	}

	return names
}

// Removes the characters in a string until and including its first underscore
func parseTableName(name string) string {
	result := ""
	runeName := []rune(name)
	foundUnderscore := false

	for i := 0; i < len(runeName); i++ {
		currentChar := string(runeName[i])
		if !foundUnderscore && currentChar == "_" {
			foundUnderscore = true
		} else if foundUnderscore {
			result += currentChar
		}
	}
	return result
}

func createTableGeneric(db *sql.DB, table_name string) error {

	query := "CREATE TABLE IF NOT EXISTS " + table_name + " (Id int unsigned NOT NULL AUTO_INCREMENT, SearchString VARCHAR(255), Site VARCHAR(255), URL VARCHAR(2040), pictureURL VARCHAR(2040), Title VARCHAR(2040), Description VARCHAR(6120), Price float, Category VARCHAR(255), PRIMARY KEY(Id))"

	create, err := db.Prepare(query)
	defer create.Close()
	if err != nil {
		log.Fatal("Creation of table went wrong in preperation stage:", err)
	}

	_, err = create.Exec()
	if err != nil {
		log.Fatal("Creation of table went wrong in execution stage:", err)
	}

	return nil
}

func createTableBlocket(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Blocket (
                          Id int unsigned NOT NULL AUTO_INCREMENT, 
                          Title varchar(30), 
                          Picture varchar(2048),
                          Description text,
                          Price int(11),
                          SellerLocation varchar(2048),
                          DateAdded varchar(30),
                          Category varchar(30),
                          SellerAccount varchar(2048),
                          AdURL varchar(2048),
                          PRIMARY KEY (Id)
                  )`

	create, err := db.Prepare(query)
	defer create.Close()
	if err != nil {
		log.Fatal("Creation of blocket table went wrong in preperation stage: %w", err)
	}

	_, err = create.Exec()
	if err != nil {
		log.Fatal("Creation of blocket table went wrong in execution stage: %w", err)
	}

	return nil
}

func createTableEbay(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Ebay (
                          Id int unsigned NOT NULL AUTO_INCREMENT, 
                          Title varchar(30), 
                          Picture varchar(2048),
                          Description text,
                          Price int(11),
                          Quality varchar(30),
                          SellerLocation varchar(2048),
                          TimeLeft int(11),
                          Category varchar(30),
                          SellerAccount varchar(2048),
                          Views int(11),
                          Quantity int(11),
                          Returns boolean,
                          Delivery boolean,
                          AdURL varchar(2048),
                          PRIMARY KEY (Id)
                  )`

	create, err := db.Prepare(query)
	defer create.Close()
	if err != nil {
		log.Fatal("Creation of Ebay table went wrong in preperation stage: %w", err)
	}

	_, err = create.Exec()
	if err != nil {
		log.Fatal("Creation of Ebay table went wrong in execution stage: %w", err)
	}

	return nil
}

func createTableFacebook(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Facebook (
                          Id int unsigned NOT NULL AUTO_INCREMENT, 
                          Title varchar(30), 
                          Picture varchar(2048),
                          Description text,
                          Price int(11),
                          SellerLocation varchar(2048),
                          FacebookUsername varchar(100),
                          Quality varchar(30),
                          SellerAccount varchar(2048),
                          AdURL varchar(2048),
                          PRIMARY KEY (Id)
                  )`

	create, err := db.Prepare(query)
	defer create.Close()
	if err != nil {
		log.Fatal("Creation of Facebook table went wrong in preperation stage: %w", err)
	}

	_, err = create.Exec()
	if err != nil {
		log.Fatal("Creation of Facebook table went wrong in execution stage: %w", err)
	}

	return nil
}

func createTableTradera(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Tradera (
                          Id int unsigned NOT NULL AUTO_INCREMENT, 
                          Title varchar(30), 
                          Picture varchar(2048),
                          Description text,
                          Price int(11),
                          SellerLocation varchar(2048),
                          SellerAccount varchar(2048),
                          Delivery boolean,
                          TimeLeft int(11),
                          DateAdded varchar(30),
                          Views int(11),
                          Category varchar(30),
                          Quality varchar(30),
                          AdURL varchar(2048),
                          PRIMARY KEY (Id)
                  )`

	create, err := db.Prepare(query)
	defer create.Close()
	if err != nil {
		log.Fatal("Creation of Tradera table went wrong in preperation stage: %w", err)
	}

	_, err = create.Exec()
	if err != nil {
		log.Fatal("Creation of Tradera table went wrong in execution stage: %w", err)
	}

	return nil
}
