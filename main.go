package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
	"bufio"

    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    Id                 int
    Title, Description string
    Quantity           int
    Price              float64
    Active             bool
}

func main() {
	db, err := sql.Open("mysql", "root:@/examGo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	actions := map[int]func(){
		1: addProduct(db),
		2: displayProducts(db),
		3: deactivateProduct(db),
		7: func() { os.Exit(0) },
	}

	for {
		fmt.Println("1- Add a product\n2- Display all products\n3- Delete a product\n7- Quit")
		var choice int
		fmt.Scan(&choice)
		if action, ok := actions[choice]; ok {
			action()
		} else {
			fmt.Println("Invalid choice. Please choose a valid option.")
		}
	}
}


// Reads product details from the database and displays them
func displayProducts(db *sql.DB) func() {
    return func() {
        rows, err := db.Query("SELECT id, title, description, quantity, price, active FROM Products")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        for rows.Next() {
            var product Product
            err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Quantity, &product.Price, &product.Active)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("Id: %d\nTitle: %s\nDescription: %s\nQuantity: %d\nPrice: %.2f\nActive: %t\n\n", product.Id, product.Title, product.Description, product.Quantity, product.Price, product.Active)
        }
    }
}

// Prompts the user to enter details about a product and adds it to the database
func addProduct(db *sql.DB) func() {
    return func() {
        product := Product{}
        scanner := bufio.NewScanner(os.Stdin)

        fmt.Print("Enter product title: ")
        scanner.Scan()
        product.Title = scanner.Text()

        fmt.Print("Enter product description: ")
        scanner.Scan()
        product.Description = scanner.Text()

        fmt.Print("Enter product quantity: ")
        fmt.Scan(&product.Quantity)

        fmt.Print("Enter product price: ")
        fmt.Scan(&product.Price)

        _, err := db.Exec(`INSERT INTO Products (title, description, quantity, price, active) VALUES (?, ?, ?, ?, ?)`,
            product.Title, product.Description, product.Quantity, product.Price, true)

        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Product inserted successfully")
    }
}

// Prompts the user to enter an ID and deactivates the corresponding product in the database
func deactivateProduct(db *sql.DB) func() {
    return func() {
        var id int
        fmt.Print("Enter the ID of the product you want to deactivate: ")
        fmt.Scan(&id)

        _, err := db.Exec("UPDATE Products SET active = ? WHERE id = ?", false, id)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Product deactivated successfully")
    }
}