package models

import (
	"database/sql"
)

type Order struct {
	Id           int
	CustomerId   int
	ProductId    int
	Quantity     int
	Price        float64
	PurchaseDate string
}

func (o *Order) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.Id, &o.CustomerId, &o.ProductId, &o.Quantity, &o.Price, &o.PurchaseDate)
}

func (o *Order) TableName() string {
	return "Orders"
}
