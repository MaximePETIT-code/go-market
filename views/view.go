package views

import (
	"bufio"
	"fmt"
)

func DisplayMenu() {
	fmt.Println("1- Add a product")
	fmt.Println("2- Display all products")
	fmt.Println("3- Modify a product")
	fmt.Println("4- Deactivate a product")
	fmt.Println("5- Add a customer")
	fmt.Println("6- Display all customers")
	fmt.Println("7- Modify a customer")
	fmt.Println("8- Export all products to CSV")
	fmt.Println("9- Export all customers to CSV")
	fmt.Println("10- Add an order")
	fmt.Println("11- Export all orders to CSV")
	fmt.Println("12- Quit")
}

func DisplayMessage(message string) {
	fmt.Println(message)
}

func PromptMessage(message string) {
	fmt.Print(message)
}

func PromptForFieldUpdate(scanner *bufio.Scanner, fields map[string]string) (string, string) {
	fmt.Println("Which field do you want to update?")
	for k, v := range fields {
		fmt.Printf("%s- %s\n", k, v)
	}
	scanner.Scan()
	choice := scanner.Text()
	field, ok := fields[choice]
	if !ok {
		return "", ""
	}
	fmt.Printf("Enter new %s: ", field)
	scanner.Scan()
	value := scanner.Text()
	return field, value
}
