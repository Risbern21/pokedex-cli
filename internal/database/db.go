package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Client() *gorm.DB {
	return db
}

func Connect() {
	var err error

	myLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err = gorm.Open(sqlite.Open("pokeballs.db"), &gorm.Config{
		Logger: myLogger,
	})
	if err != nil {
		log.Fatalf("unable to open database : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error while opening database : %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("error while pinging database : %v", err)
	}
	log.Println("successfully connected to database")
}
