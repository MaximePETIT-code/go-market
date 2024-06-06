package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"market/views"
)

type Customer struct {
	Id                                         int
	FirstName, LastName, Phone, Address, Email string
}

func (c *Customer) Scan(rows *sql.Rows) error {
	return rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Address, &c.Email)
}

func (c *Customer) Display() {
	fmt.Printf("ID: %d\nFirst Name: %s\nLast Name: %s\nPhone: %s\nAddress: %s\nEmail: %s\n\n", c.Id, c.FirstName, c.LastName, c.Phone, c.Address, c.Email)
}

func (c *Customer) TableName() string {
	return "Customers"
}

func (c *Customer) PromptForValues() {
	scanner := bufio.NewScanner(os.Stdin)
	views.PromptMessage("Enter customer first name: ")
	scanner.Scan()
	c.FirstName = scanner.Text()

	views.PromptMessage("Enter customer last name: ")
	scanner.Scan()
	c.LastName = scanner.Text()

	views.PromptMessage("Enter customer phone: ")
	scanner.Scan()
	c.Phone = scanner.Text()

	views.PromptMessage("Enter customer address: ")
	scanner.Scan()
	c.Address = scanner.Text()

	views.PromptMessage("Enter customer email: ")
	scanner.Scan()
	c.Email = scanner.Text()
}

func (c *Customer) PromptForUpdate() (string, interface{}) {
	scanner := bufio.NewScanner(os.Stdin)
	field, value := views.PromptForFieldUpdate(scanner, map[string]string{
		"1": "firstName",
		"2": "lastName",
		"3": "phone",
		"4": "address",
		"5": "email",
	})
	return field, value
}
