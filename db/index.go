package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("Failed to connect to database")
		panic("Failed to connect to database") //panic
	}

	// Enable uui-ossp extension
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		fmt.Println("can't install uuid extension")
		panic(err)
	}

	err = DBConn.AutoMigrate(&User{}, &SearchSetting{}, &CrawledUrl{})
	// AutoMigrate creates tables based on the struct/model
	// &User{} is a pointer to the User struct/model
	// &SearchSetting{} is a pointer to the SearchSetting struct/model
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConn
}
