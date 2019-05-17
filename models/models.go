package models

import (
	"zlp-demo-golang/provider"
	"log"

	"github.com/jinzhu/gorm"
)

func InitModels() {
	err := provider.UseDB(func(db *gorm.DB) {
		// db.DropTableIfExists(&Order{})
		db.AutoMigrate(&Order{})
	})

	if err != nil {
		log.Fatal(err)
	}
}
