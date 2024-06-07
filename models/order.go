package models

import (
	"database/sql"
	"fmt"

	"github.com/jung-kurt/gofpdf"
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

func (o *Order) GeneratePDFContent(pdf *gofpdf.Fpdf) {
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Customer Details:")
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("First Name: %s", o.Customer.FirstName))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Last Name: %s", o.Customer.LastName))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Email: %s", o.Customer.Email))
    pdf.Ln(10)

    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Order Details:")
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("Product Name: %s", o.Product.Title))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Product ID: %d", o.Product.Id))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Quantity: %d", o.Quantity))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Price: %.2f", o.Price))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Total Price: %.2f", o.Price*float64(o.Quantity)))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Purchase Date: %s", o.PurchaseDate))
    pdf.Ln(10)
}