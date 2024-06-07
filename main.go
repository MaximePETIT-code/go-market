package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"market/controllers"
	"market/models"
	"market/views"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@/examGo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	actions := map[int]func(){
		1:  func() { controllers.HandleAction(controllers.AddEntity, &models.Product{}, db) },
		2:  func() { controllers.HandleAction(controllers.DisplayEntities, &models.Product{}, db) },
		3:  func() { controllers.HandleAction(controllers.ModifyEntity, &models.Product{}, db) },
		4:  func() { controllers.HandleAction(controllers.DeactivateEntity, &models.Product{}, db) },
		5:  func() { controllers.HandleAction(controllers.AddEntity, &models.Customer{}, db) },
		6:  func() { controllers.HandleAction(controllers.DisplayEntities, &models.Customer{}, db) },
		7:  func() { controllers.HandleAction(controllers.ModifyEntity, &models.Customer{}, db) },
		8:  func() { controllers.ExportToCSV(&models.Product{}, db, "exports/products.csv") },
		9:  func() { controllers.ExportToCSV(&models.Customer{}, db, "exports/customers.csv") },
		10: func() { controllers.AddOrder(db) },
		11: func() { controllers.ExportToCSV(&models.Order{}, db, "exports/orders.csv") },
		12: func() { os.Exit(0) },
	}

	for {
		views.DisplayMenu()
		var choice int
		fmt.Scan(&choice)
		if action, ok := actions[choice]; ok {
			action()
		} else {
			views.DisplayMessage("Invalid choice. Please choose a valid option.")
		}
	}
}
