package respository

import (
	"zlp-demo-golang/models"
	"zlp-demo-golang/provider"

	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
)

const (
	OrderPerPage = 10
)

type orderRespository struct{}

// OrderRespository ...
var OrderRespository = &orderRespository{}

// SaveOrder save order to database
func (orderRespo *orderRespository) SaveOrder(data string) {
	provider.UseDB(func(db *gorm.DB) {
		embeddata := gjson.Get(data, "embeddata").String()
		db.Save(&models.Order{
			Apptransid:  gjson.Get(data, "apptransid").String(),
			Zptransid:   gjson.Get(data, "zptransid").String(),
			Description: gjson.Get(embeddata, "description").String(), // Nhúng description vào embeddata vì trong callback data không có trường này
			Amount:      gjson.Get(data, "amount").Int(),
			Timestamp:   gjson.Get(data, "servertime").Int(),
			Channel:     gjson.Get(data, "channel").Int(),
		})
	})
}

type PaginationResult struct {
	CurrentPage  int            `json:"currentPage"`
	TotalOrder   int            `json:"totalOrder"`
	OrderPerPage int            `json:"orderPerPage"`
	Orders       []models.Order `json:"orders"`
}

// Paginate with page >= 0 (zero-based)
func (orderRespo *orderRespository) Paginate(page int) PaginationResult {
	db, _ := provider.NewDB()
	orders := make([]models.Order, 0)
	totalOrder := 0

	db.Offset((page - 1) * OrderPerPage).Limit(OrderPerPage).Order("timestamp desc").Find(&orders)
	db.Model(&models.Order{}).Count(&totalOrder)

	return PaginationResult{
		CurrentPage:  page,
		TotalOrder:   totalOrder,
		Orders:       orders,
		OrderPerPage: OrderPerPage,
	}
}
