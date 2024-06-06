package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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
	Id                                int
	FirstName, LastName, Phone, Address, Email string
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

func main() {
	db, err := sql.Open("mysql", "root:@/examGo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	actions := map[int]func(){
		1: func() { handleAction(addEntity, &Product{}, db) },
		2: func() { handleAction(displayEntities, &Product{}, db) },
		3: func() { handleAction(modifyEntity, &Product{}, db) },
		4: func() { handleAction(deactivateEntity, &Product{}, db) },
		5: func() { handleAction(addEntity, &Customer{}, db) },
		6: func() { handleAction(displayEntities, &Customer{}, db) },
		7: func() { handleAction(modifyEntity, &Customer{}, db) },
		8: func() { os.Exit(0) },
	}

	for {
		fmt.Println("1- Add a product\n2- Display all products\n3- Modify a product\n4- Deactivate a product\n5- Add a customer\n6- Display all customers\n7- Modify a customer\n8- Quit")
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
