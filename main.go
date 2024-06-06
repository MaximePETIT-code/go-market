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
}