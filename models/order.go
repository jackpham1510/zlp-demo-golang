package models

type Order struct {
	Apptransid  string `gorm:"primary_key"`
	Zptransid   string
	Description string
	Amount      int64
	Timestamp   int64 // unix timestamp
	Channel     int64
}
