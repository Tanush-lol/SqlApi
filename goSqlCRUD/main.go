package main

import(	
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	)

func renderTable(tableData [][]any){

for _, row := range tableData {
    for _, val := range row {
        switch v := val.(type) {
        case []byte:
            fmt.Printf("%-25s", string(v))
        case time.Time:
            fmt.Printf("%-12s", v.Format("2006-01-02")) 
        default:
            fmt.Printf("%-12v", v)
        }
    }
    fmt.Println()
}
	fmt.Println("\n")
}

func main(){

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
	db, err := sql.Open("postgres", connStr)

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

	tx, err := db.Begin() //opens up db for editing stuff
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	
	fmt.Println("Here have a look at the database")
	rows,err := tx.Query(`
		SELECT * FROM employee
		`)
	if err != nil{
		tx.Rollback()
		log.Fatal("query failed",err)
	}
	defer rows.Close()

	cols,err :=rows.Columns();
	fmt.Println("The columns are",cols)

	var tableData [][]any

	for rows.Next(){
		rowValues :=make([]any,len(cols))

		rowPointers:= make([]any,len(cols))
		
		for i :=range rowValues{
			rowPointers[i]=&rowValues[i]
		}

	err:=rows.Scan(rowPointers...)
		if err!=nil{
			fmt.Println("Error ",err)
			break
		}
	tableData= append(tableData,rowValues)
	}
	renderTable(tableData)

	fmt.Println("Here's the menu for CRUD operations on the sql database")
}

