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

type Customer struct {
    Id        int
    FirstName string
    LastName  string
    Phone     string
    Address   string
    Email     string
}

func main() {
    db, err := sql.Open("mysql", "root:@/examGo")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    actions := map[int]func(){
        1: addProduct(db),
        2: displayProducts(db),
		3: modifyProduct(db),
        4: deactivateProduct(db),
        5: addCustomer(db),
        6: displayCustomers(db),
        7: modifyCustomer(db),
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

func modifyProduct(db *sql.DB) func() {
    return func() {
        var id int
        fmt.Print("Enter the ID of the product you want to modify: ")
        fmt.Scan(&id)

        fmt.Println("Which field do you want to update?")
        fmt.Println("1- Title\n2- Description\n3- Quantity\n4- Price\n5- Active")
        var choice int
        fmt.Scan(&choice)

        scanner := bufio.NewScanner(os.Stdin)
        var field, value string

        switch choice {
        case 1:
            field = "title"
            fmt.Print("Enter new title: ")
        case 2:
            field = "description"
            fmt.Print("Enter new description: ")
        case 3:
            field = "quantity"
            fmt.Print("Enter new quantity: ")
            scanner.Scan()
            value = scanner.Text()
        case 4:
            field = "price"
            fmt.Print("Enter new price: ")
            scanner.Scan()
            value = scanner.Text()
        case 5:
            field = "active"
            fmt.Print("Enter new active status (true/false): ")
        default:
            fmt.Println("Invalid choice. Please choose a valid option.")
            return
        }

        scanner.Scan()
        value = scanner.Text()

        query := fmt.Sprintf("UPDATE Products SET %s = ? WHERE id = ?", field)
        _, err := db.Exec(query, value, id)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Product modified successfully")
    }
}

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

func addCustomer(db *sql.DB) func() {
    return func() {
        customer := Customer{}
        scanner := bufio.NewScanner(os.Stdin)

        fmt.Print("Enter customer first name: ")
        scanner.Scan()
        customer.FirstName = scanner.Text()

        fmt.Print("Enter customer last name: ")
        scanner.Scan()
        customer.LastName = scanner.Text()

        fmt.Print("Enter customer phone: ")
        scanner.Scan()
        customer.Phone = scanner.Text()

        fmt.Print("Enter customer address: ")
        scanner.Scan()
        customer.Address = scanner.Text()

        fmt.Print("Enter customer email: ")
        scanner.Scan()
        customer.Email = scanner.Text()

        _, err := db.Exec(`INSERT INTO Customers (firstName, lastName, phone, address, email) VALUES (?, ?, ?, ?, ?)`,
            customer.FirstName, customer.LastName, customer.Phone, customer.Address, customer.Email)

        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Customer inserted successfully")
    }
}

func displayCustomers(db *sql.DB) func() {
    return func() {
        rows, err := db.Query("SELECT id, firstName, lastName, phone, address, email FROM Customers")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        for rows.Next() {
            var customer Customer
            err := rows.Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Phone, &customer.Address, &customer.Email)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("Id: %d\nFirst Name: %s\nLast Name: %s\nPhone: %s\nAddress: %s\nEmail: %s\n\n", customer.Id, customer.FirstName, customer.LastName, customer.Phone, customer.Address, customer.Email)
        }
    }
}

func modifyCustomer(db *sql.DB) func() {
    return func() {
        var id int
        fmt.Print("Enter the ID of the customer you want to modify: ")
        fmt.Scan(&id)

        fmt.Println("Which field do you want to update?")
        fmt.Println("1- First Name\n2- Last Name\n3- Phone\n4- Address\n5- Email")
        var choice int
        fmt.Scan(&choice)

        scanner := bufio.NewScanner(os.Stdin)
        var field, value string

        switch choice {
        case 1:
            field = "firstName"
            fmt.Print("Enter new first name: ")
        case 2:
            field = "lastName"
            fmt.Print("Enter new last name: ")
        case 3:
            field = "phone"
            fmt.Print("Enter new phone: ")
        case 4:
            field = "address"
            fmt.Print("Enter new address: ")
        case 5:
            field = "email"
            fmt.Print("Enter new email: ")
        default:
            fmt.Println("Invalid choice. Please choose a valid option.")
            return
        }

        scanner.Scan()
        value = scanner.Text()

        query := fmt.Sprintf("UPDATE Customers SET %s = ? WHERE id = ?", field)
        _, err := db.Exec(query, value, id)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Customer modified successfully")
    }
}