package models

import (
	"database/sql"
	"fmt"
)

type Order struct {
	Id           int
	Customer   Customer
	Product    Product
	Quantity     int
	Price        float64
	PurchaseDate string
}

func (o *Order) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.Id, &o.Customer.Id, &o.Product.Id, &o.Quantity, &o.Price, &o.PurchaseDate)
}

func (o *Order) TableName() string {
	return "Orders"
}

func (o *Order) Display() {
    fmt.Printf("Order ID: %d\nCustomer ID: %d\nProduct ID: %d\nQuantity: %d\n\n", o.Id, o.Customer.Id, o.Product.Id, o.Quantity)
}
