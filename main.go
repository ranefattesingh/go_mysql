package main

import (
	"fmt"
	"bufio"
	"os"
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type student struct {
	id 			int
	name		string
	age			byte
	percentage	float32
}

const CLIENT = "mysql"
const USER = "root"
const PWD = "home"
const DB = "test"


func _verifyConnecton() {
	db, err := sql.Open(CLIENT, USER + ":" + PWD + "@/" + DB)
	defer db.Close()

	if (err != nil) {
		log.Fatal("Error connecting database.")
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)

	if (err != nil) {
		log.Fatal("Error connecting database.")
	}

	fmt.Println("Mysql version is ", version)
}

func _insertStudent(newStudent *student) {
	fmt.Println((*newStudent).name)
	fmt.Println((*newStudent).age)
	fmt.Println((*newStudent).percentage)

	db, err := sql.Open(CLIENT, USER + ":" + PWD + "@/" + DB)
	defer db.Close()

	if err != nil {
		log.Fatal("Error connecting database.")
	}

	getMaxIdQuery := "SELECT MAX(Id) FROM Student"
	rows, err := db.Query(getMaxIdQuery)
	var maxId int = 0

	/* If instead of for because we expect only one row here*/
	if rows.Next() {
		rows.Scan(&maxId)
	}

	insertQuery := fmt.Sprintf("INSERT INTO Student(Id, Name, Age, Percentage) VALUES (%d, '%s', %d, %f)",
		maxId + 1, 
		(*newStudent).name, 
		(*newStudent).age, 
		(*newStudent).percentage,
	)

	res, err := db.Exec(insertQuery)

	if err != nil {
		log.Fatal("Error inserting values")
	}

	if rows, _ := res.RowsAffected(); rows == 1 {
		fmt.Println("Student data is inserted.")
	}
}

func fetchStudent(name string) []student {
	db, err := sql.Open(CLIENT, fmt.Sprintf("%s:%s@/%s", USER, PWD, DB))
	defer db.Close()

	if err != nil {
		log.Fatal("Error connecting database!")
	}

	fetchQuery := fmt.Sprintf("SELECT Id, Name, Age, Percentage FROM Student WHERE Name = '%s'", name) 
	fmt.Println(fetchQuery)
	
	var students []student
	rows, err := db.Query(fetchQuery)
	for rows.Next() {
		var s student
		fmt.Scan(&s.id, &s.name, &s.age, &s.percentage)
		students = append(students, s)
	}

	fmt.Println(students)
	return students
}

func main() {
	// _verifyConnecton()
	// var newStudent student
	// fmt.Print("Enter student name, age and percentage respectively: ")
	// fmt.Scan(&newStudent.name, &newStudent.age, &newStudent.percentage)
	// insertStudent(&newStudent)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter student name: ")
	scanner.Scan()
	name := scanner.Text()
	fetchStudent(name)
}

