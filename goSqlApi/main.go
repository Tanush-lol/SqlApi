package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type Employee struct {
	ID     int
	Name   string
	Email  string
	Salary int
}

func main() {
	// Get database credentials from user
	reader := bufio.NewReader(os.Stdin) //reads from input and stores it in var reader

	fmt.Print("Enter database name: ")
	dbname, _ := reader.ReadString('\n')
	dbname = strings.TrimSpace(dbname)

	fmt.Print("Enter username: ")
	user, _ := reader.ReadString('\n') //takes username for db
	user = strings.TrimSpace(user)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n') //takes user password for db
	password = strings.TrimSpace(password)

	// Build connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost sslmode=disable",
		user, password, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connStr)//connection to the db returns db or err so we write it like 
																					//like this
	//in pq library the init function demands an initial driver name which in case of postgres is postgres
	//like while opening postgressql we write psql -U postgres this it the initial driver it demands to 
	//open the psql entry point

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	fmt.Println("Connected to database successfully!")

	// Start transaction
	tx, err := db.Begin() //opens up db for editing stuff
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	// Operation 1: Query all employees ordered by name
	fmt.Println("\n=== All Employees (Ordered by Name) ===")
	rows1, err := tx.Query(`
		SELECT id, name, email, salary 
		FROM employee 
		ORDER BY name
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal("Query 1 failed:", err)
	}
	defer rows1.Close()

	for rows1.Next() {
		var emp Employee
		err := rows1.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Salary)
		if err != nil {
			tx.Rollback()
			log.Fatal("Failed to scan row:", err)
		}
		fmt.Printf("ID: %d | Name: %-20s | Email: %-30s | Salary: $%d\n",
			emp.ID, emp.Name, emp.Email, emp.Salary)
	}

	// Operation 2: Stream employees (simulate with another query)
	fmt.Println("\n=== All Employees (Stream-like) ===")
	rows2, err := tx.Query("SELECT name, salary FROM employee")
	if err != nil {
		tx.Rollback()
		log.Fatal("Query 2 failed:", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var name string
		var salary int
		err := rows2.Scan(&name, &salary)
		if err != nil {
			tx.Rollback()
			log.Fatal("Failed to scan row:", err)
		}
		fmt.Printf("%s earns $%d\n", name, salary)
	}

	// Operation 3: Double all salaries
	fmt.Println("\n=== Doubling all employees' salaries ===")
	result, err := tx.Exec("UPDATE employee SET salary = salary * 2")
	if err != nil {
		tx.Rollback()
		log.Fatal("Update failed:", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Updated %d records\n", rowsAffected)

	// Operation 4: Query specific employee's salary (using Tanush Roy as example)
	fmt.Println("\n=== Querying specific employee salary ===")
	var mySalary int
	err = tx.QueryRow(`
		SELECT salary 
		FROM employee 
		WHERE name = 'Tanush Roy'
	`).Scan(&mySalary)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Employee 'Tanush Roy' not found")
		} else {
			tx.Rollback()
			log.Fatal("Query failed:", err)
		}
	} else {
		fmt.Printf("Tanush Roy now earns $%d\n", mySalary)
	}

	// Operation 5: Get top earner
	fmt.Println("\n=== Top Earner ===")
	var topName string
	var topSalary int
	err = tx.QueryRow(`
		SELECT name, salary 
		FROM employee 
		WHERE salary = (SELECT MAX(salary) FROM employee) 
		LIMIT 1
	`).Scan(&topName, &topSalary)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No employees found")
		} else {
			tx.Rollback()
			log.Fatal("Top earner query failed:", err)
		}
	} else {
		fmt.Printf("Top earner is %s with a salary of $%d\n", topName, topSalary)
	}

	fmt.Println("\n=== Employee Table Columns ===")
	rows3, err := tx.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'employee' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal("Column query failed:", err)
	}
	defer rows3.Close()

	for rows3.Next() {
		var colName, dataType string
		err := rows3.Scan(&colName, &dataType)
		if err != nil {
			tx.Rollback()
			log.Fatal("Failed to scan column info:", err)
		}
		fmt.Printf("%-15s : %s\n", colName, dataType)
	}

	// Operation 7: Get all employees after update
	fmt.Println("\n=== All Employees After Salary Update ===")
	rows4, err := tx.Query(`
		SELECT id, name, email, salary 
		FROM employee 
		ORDER BY salary DESC
	`)
	if err != nil {
		tx.Rollback()
		log.Fatal("Final query failed:", err)
	}
	defer rows4.Close()

	for rows4.Next() {
		var emp Employee
		err := rows4.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Salary)
		if err != nil {
			tx.Rollback()
			log.Fatal("Failed to scan final row:", err)
		}
		fmt.Printf("ID: %2d | Name: %-20s | Salary: $%d\n",
			emp.ID, emp.Name, emp.Salary)
	}

	fmt.Print("\nDo you want to commit these changes? (yes/no): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "yes" || answer == "y" {
		err = tx.Commit()
		if err != nil {
			log.Fatal("Failed to commit transaction:", err)
		}
		fmt.Println("Transaction committed successfully!")
	} else {
		// Rollback transaction
		err = tx.Rollback()
		if err != nil {
			log.Fatal("Failed to rollback transaction:", err)
		}
		fmt.Println("Transaction rolled back. No changes made.")
	}
}
