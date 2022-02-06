package database

import "gorm.io/gorm"

func Exists(tx *gorm.DB, query interface{}, args ...interface{}) bool {
	var exists bool
	tx.Where(query, args...).Select("count(*) > 0").Find(&exists)
	return exists
}
