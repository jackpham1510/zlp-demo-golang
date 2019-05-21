package models

type Order struct {
	Apptransid  string `gorm:"primary_key" json:"apptransid"`
	Zptransid   string `json:"zptransid"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Timestamp   int64  `json:"timestamp"`
	Channel     int64  `json:"channel"`
}
