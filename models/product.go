package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"market/views"
	"os"
	"strconv"
)

type Product struct {
	Id                 int
	Title, Description string
	Quantity           int
	Price              float32
	Active             bool
}

func (p *Product) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Id, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Active)
}

func (p *Product) Display() {
	fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\nQuantity: %d\nPrice: %.2f\nActive: %v\n\n", p.Id, p.Title, p.Description, p.Quantity, p.Price, p.Active)
}

func (p *Product) TableName() string {
	return "Products"
}

func (p *Product) PromptForValues() {
	scanner := bufio.NewScanner(os.Stdin)
	views.PromptMessage("Enter product title: ")
	scanner.Scan()
	p.Title = scanner.Text()

	views.PromptMessage("Enter product description: ")
	scanner.Scan()
	p.Description = scanner.Text()

	views.PromptMessage("Enter product quantity: ")
	fmt.Scan(&p.Quantity)

	views.PromptMessage("Enter product price: ")
	fmt.Scan(&p.Price)

	p.Active = true
}

func (p *Product) PromptForUpdate() (string, interface{}) {
	scanner := bufio.NewScanner(os.Stdin)
	field, value := views.PromptForFieldUpdate(scanner, map[string]string{
		"1": "title",
		"2": "description",
		"3": "quantity",
		"4": "price",
		"5": "active",
	})

	switch field {
	case "quantity":
		v, _ := strconv.Atoi(value)
		return field, v
	case "price":
		v, _ := strconv.ParseFloat(value, 32)
		return field, float32(v)
	case "active":
		v, _ := strconv.ParseBool(value)
		return field, v
	default:
		return field, value
	}
}
