package controllers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"market/models"
	"market/views"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Entity interface {
	Scan(*sql.Rows) error
	Display()
	PromptForValues()
	TableName() string
	PromptForUpdate() (string, interface{})
}

func HandleAction(action func(Entity, *sql.DB), entity Entity, db *sql.DB) {
	action(entity, db)
}

func DisplayEntities(entity Entity, db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", entity.TableName()))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := entity.Scan(rows)
		if err != nil {
			log.Fatal(err)
		}
		entity.Display()
	}
}

func AddEntity(entity Entity, db *sql.DB) {
	entity.PromptForValues()
	columns, values := getColumnsAndValues(entity)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, entity.TableName(), columns, values)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func ModifyEntity(entity Entity, db *sql.DB) {
	var id int
	views.PromptMessage("Enter the ID of the entity you want to modify: ")
	fmt.Scan(&id)

	field, value := entity.PromptForUpdate()
	if field != "" {
		_, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", entity.TableName(), field), value, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		views.DisplayMessage("Invalid field choice.")
	}
}

func DeactivateEntity(entity Entity, db *sql.DB) {
	var id int
	views.PromptMessage("Enter the ID of the entity you want to deactivate: ")
	fmt.Scan(&id)

	_, err := db.Exec(fmt.Sprintf("UPDATE %s SET active = 0 WHERE id = ?", entity.TableName()), id)
	if err != nil {
		log.Fatal(err)
	}
}

func ExportToCSV(entity Entity, db *sql.DB, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers to CSV
	var headers []string
	switch entity.(type) {
	case *models.Product:
		headers = []string{"Id", "Title", "Description", "Quantity", "Price", "Active"}
	case *models.Customer:
		headers = []string{"Id", "FirstName", "LastName", "Phone", "Address", "Email"}
	}
	writer.Write(headers)

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", entity.TableName()))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := entity.Scan(rows)
		if err != nil {
			log.Fatal(err)
		}
		var record []string
		switch e := entity.(type) {
		case *models.Product:
			record = []string{strconv.Itoa(e.Id), e.Title, e.Description, strconv.Itoa(e.Quantity), fmt.Sprintf("%.2f", e.Price), strconv.FormatBool(e.Active)}
		case *models.Customer:
			record = []string{strconv.Itoa(e.Id), e.FirstName, e.LastName, e.Phone, e.Address, e.Email}
		}
		writer.Write(record)
	}
}

func AddOrder(db *sql.DB) {
	order := &models.Order{}
	fmt.Print("Enter customer ID: ")
	fmt.Scan(&order.CustomerId)
	fmt.Print("Enter product ID: ")
	fmt.Scan(&order.ProductId)
	fmt.Print("Enter quantity: ")
	fmt.Scan(&order.Quantity)

	// Get the product details from the database
	var productName string
	var availableQuantity int
	row := db.QueryRow("SELECT title, price, quantity FROM Products WHERE id = ?", order.ProductId)
	err := row.Scan(&productName, &order.Price, &availableQuantity)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the available quantity is sufficient
	if availableQuantity < order.Quantity {
		fmt.Println("Insufficient quantity available. Order cannot be placed.")
		return
	}

	// Update the quantity in the database
	_, err = db.Exec("UPDATE Products SET quantity = ? WHERE id = ?", availableQuantity-order.Quantity, order.ProductId)
	if err != nil {
		log.Fatal(err)
	}

	// Get the current date
	order.PurchaseDate = time.Now().Format("2006-01-02")

	columns := "customerId, productId, quantity, price, purchaseDate"
	values := fmt.Sprintf("%d, %d, %d, %.2f, '%s'", order.CustomerId, order.ProductId, order.Quantity, order.Price, order.PurchaseDate)
	query := fmt.Sprintf(`INSERT INTO Orders (%s) VALUES (%s)`, columns, values)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Get customer details from the database
	var firstName, lastName, customerEmail string
	row = db.QueryRow("SELECT firstName, lastName, email FROM Customers WHERE id = ?", order.CustomerId)
	err = row.Scan(&firstName, &lastName, &customerEmail)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add customer details to the PDF
	pdf.Cell(40, 10, "Customer Details:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("First Name: %s", firstName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Last Name: %s", lastName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Email: %s", customerEmail))
	pdf.Ln(10)

	// Add order details to the PDF
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Order Details:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Product Name: %s", productName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Product ID: %d", order.ProductId))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Quantity: %d", order.Quantity))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Price: %.2f", order.Price))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Total Price: %.2f", order.Price*float64(order.Quantity)))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Purchase Date: %s", order.PurchaseDate))
	pdf.Ln(10)

	// Save the PDF to a file
	filename := fmt.Sprintf("exports/orders/order_%s.pdf", time.Now().Format("20060102_150405"))
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func getColumnsAndValues(entity Entity) (string, string) {
	switch e := entity.(type) {
	case *models.Product:
		return "title, description, quantity, price, active", fmt.Sprintf("'%s', '%s', %d, %.2f, %v", e.Title, e.Description, e.Quantity, e.Price, e.Active)
	case *models.Customer:
		return "firstName, lastName, phone, address, email", fmt.Sprintf("'%s', '%s', '%s', '%s', '%s'", e.FirstName, e.LastName, e.Phone, e.Address, e.Email)
	}
	return "", ""
}