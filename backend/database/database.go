package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var globDatabase *gorm.DB

func newDatabaseConnection(dialector gorm.Dialector) (*gorm.DB, error) {
	return gorm.Open(dialector)
}

func NewSqlite3Connection(databaseFile string) (*gorm.DB, error) {
	return newDatabaseConnection(sqlite.Open(databaseFile))
}

func GetDB() *gorm.DB {
	return globDatabase
}

func SetGlobDB(database *gorm.DB) {
	globDatabase = database
}
