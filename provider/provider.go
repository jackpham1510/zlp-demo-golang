package provider

import (
	"esale-issues/config"
	"log"

	"github.com/jinzhu/gorm"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.String("db.connstring"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

func UseDB(fn func(*gorm.DB)) error {
	db, err := NewDB()
	if err != nil {
		return err
	}
	defer db.Close()
	fn(db)
	return nil
}
