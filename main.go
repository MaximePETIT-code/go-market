package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jung-kurt/gofpdf"
)

type Entity interface {
	Scan(*sql.Rows) error
	Display()
	PromptForValues(*bufio.Scanner)
	TableName() string
	PromptForUpdate(*bufio.Scanner) (string, interface{})
}

type Product struct {
	Id                 int
	Title, Description string
	Quantity           int
	Price              float32
	Active             bool
}

type Customer struct {
	Id                                         int
	FirstName, LastName, Phone, Address, Email string
}

type Order struct {
	Id           int
	CustomerId   int
	ProductId    int
	Quantity     int
	Price        float64
	PurchaseDate string
}

func (p *Product) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Id, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Active)
}

func (c *Customer) Scan(rows *sql.Rows) error {
	return rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Address, &c.Email)
}

func (p *Product) Display() {
	fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\nQuantity: %d\nPrice: %.2f\nActive: %v\n\n", p.Id, p.Title, p.Description, p.Quantity, p.Price, p.Active)
}

func (c *Customer) Display() {
	fmt.Printf("ID: %d\nFirst Name: %s\nLast Name: %s\nPhone: %s\nAddress: %s\nEmail: %s\n\n", c.Id, c.FirstName, c.LastName, c.Phone, c.Address, c.Email)
}

func (p *Product) TableName() string {
	return "Products"
}

func (c *Customer) TableName() string {
	return "Customers"
}

// func (o *Order) TableName() string {
// 	return "Orders"
// }

func main() {
	db, err := sql.Open("mysql", "root:@/examGo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	actions := map[int]func(){
		1:  func() { handleAction(addEntity, &Product{}, db) },
		2:  func() { handleAction(displayEntities, &Product{}, db) },
		3:  func() { handleAction(modifyEntity, &Product{}, db) },
		4:  func() { handleAction(deactivateEntity, &Product{}, db) },
		5:  func() { handleAction(addEntity, &Customer{}, db) },
		6:  func() { handleAction(displayEntities, &Customer{}, db) },
		7:  func() { handleAction(modifyEntity, &Customer{}, db) },
		8:  func() { exportToCSV(&Product{}, db, "exports/products.csv") },
		9:  func() { exportToCSV(&Customer{}, db, "exports/customers.csv") },
		10: func() { addOrder(db) },
		11: func() { os.Exit(0) }}

	for {
		fmt.Println("1- Add a product\n2- Display all products\n3- Modify a product\n4- Deactivate a product\n5- Add a customer\n6- Display all customers\n7- Modify a customer\n8- Export all products to CSV\n9- Export all customers to CSV\n10- Add an order\n11- Quit")
		var choice int
		fmt.Scan(&choice)
		if action, ok := actions[choice]; ok {
			action()
		} else {
			fmt.Println("Invalid choice. Please choose a valid option.")
		}
	}
}

func handleAction(action func(Entity, *sql.DB), entity Entity, db *sql.DB) {
	action(entity, db)
}

func displayEntities(entity Entity, db *sql.DB) {
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

func addEntity(entity Entity, db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	entity.PromptForValues(scanner)
	columns, values := getColumnsAndValues(entity)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, entity.TableName(), columns, values)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func modifyEntity(entity Entity, db *sql.DB) {
	var id int
	fmt.Print("Enter the ID of the entity you want to modify: ")
	fmt.Scan(&id)

	scanner := bufio.NewScanner(os.Stdin)
	field, value := entity.PromptForUpdate(scanner)
	if field != "" {
		_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE id = ?`, entity.TableName(), field), value, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Invalid choice. Please choose a valid option.")
	}
}

func deactivateEntity(entity Entity, db *sql.DB) {
	var id int
	fmt.Print("Enter the ID of the entity you want to deactivate: ")
	fmt.Scan(&id)

	_, err := db.Exec(fmt.Sprintf(`UPDATE %s SET active = ? WHERE id = ?`, entity.TableName()), false, id)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Product) PromptForValues(scanner *bufio.Scanner) {
	prompt(scanner, "Enter product title: ", &p.Title)
	prompt(scanner, "Enter product description: ", &p.Description)
	fmt.Print("Enter product quantity: ")
	fmt.Scan(&p.Quantity)
	fmt.Print("Enter product price: ")
	fmt.Scan(&p.Price)
	p.Active = true
}

func (c *Customer) PromptForValues(scanner *bufio.Scanner) {
	prompt(scanner, "Enter customer first name: ", &c.FirstName)
	prompt(scanner, "Enter customer last name: ", &c.LastName)
	prompt(scanner, "Enter customer phone: ", &c.Phone)
	prompt(scanner, "Enter customer address: ", &c.Address)
	prompt(scanner, "Enter customer email: ", &c.Email)
}

func prompt(scanner *bufio.Scanner, message string, field *string) {
	fmt.Print(message)
	scanner.Scan()
	*field = scanner.Text()
}

func (p *Product) PromptForUpdate(scanner *bufio.Scanner) (string, interface{}) {
	return promptForFieldUpdate(scanner, map[string]string{
		"1": "title",
		"2": "description",
		"3": "quantity",
		"4": "price",
		"5": "active",
	})
}

func (c *Customer) PromptForUpdate(scanner *bufio.Scanner) (string, interface{}) {
	return promptForFieldUpdate(scanner, map[string]string{
		"1": "firstName",
		"2": "lastName",
		"3": "phone",
		"4": "address",
		"5": "email",
	})
}

func promptForFieldUpdate(scanner *bufio.Scanner, fields map[string]string) (string, interface{}) {
	fmt.Println("Which field do you want to update?")
	for k, v := range fields {
		fmt.Printf("%s- %s\n", k, v)
	}
	scanner.Scan()
	choice := scanner.Text()
	field, ok := fields[choice]
	if !ok {
		return "", nil
	}
	fmt.Printf("Enter new %s: ", field)
	scanner.Scan()
	value := scanner.Text()
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

func getColumnsAndValues(entity Entity) (string, string) {
	switch e := entity.(type) {
	case *Product:
		return "title, description, quantity, price, active", fmt.Sprintf("'%s', '%s', %d, %.2f, %v", e.Title, e.Description, e.Quantity, e.Price, e.Active)
	case *Customer:
		return "firstName, lastName, phone, address, email", fmt.Sprintf("'%s', '%s', '%s', '%s', '%s'", e.FirstName, e.LastName, e.Phone, e.Address, e.Email)
	}
	return "", ""
}

func exportToCSV(entity Entity, db *sql.DB, filePath string) {
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
	case *Product:
		headers = []string{"Id", "Title", "Description", "Quantity", "Price", "Active"}
	case *Customer:
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
		case *Product:
			record = []string{strconv.Itoa(e.Id), e.Title, e.Description, strconv.Itoa(e.Quantity), fmt.Sprintf("%.2f", e.Price), strconv.FormatBool(e.Active)}
		case *Customer:
			record = []string{strconv.Itoa(e.Id), e.FirstName, e.LastName, e.Phone, e.Address, e.Email}
		}
		writer.Write(record)
	}
}

func addOrder(db *sql.DB) {
    order := &Order{}
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
