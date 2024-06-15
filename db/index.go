package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB // DBConn is a global variable that holds the database connection.

// InitDB is a function that initializes the database connection and sets up the necessary tables.
// It retrieves the database URL from the environment variables and attempts to connect to the database.
// If the connection fails, it prints an error message and panics.
// It then enables the "uuid-ossp" extension in the database. If this fails, it prints an error message and panics.
// Finally, it attempts to auto-migrate the User, SearchSettings, CrawledUrl, and SearchIndex tables.
// If the migration fails, it prints an error message and panics.
//
// This function does not take any parameters and does not return any values.
func InitDB() {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("Failed to connect to database")
		panic("Failed to connect to database") // panic is used to stop the execution of the program
	}

	// Enable uuid-ossp extension
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		fmt.Println("Failed to enable uuid-ossp extension")
		panic(err)
	}

	err = DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledUrl{}, &SearchIndex{})
	if err != nil {
		fmt.Println("Failed to migrate")
		panic(err)
	}
}

// GetDB is a function that returns the current database connection.
// It does not take any parameters.
//
// Returns:
// *gorm.DB: The current database connection.
func GetDB() *gorm.DB {
	return DBConn
}
