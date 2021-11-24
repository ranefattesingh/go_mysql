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


func verifyConnecton() bool {
	db, err := sql.Open(CLIENT, USER + ":" + PWD + "@/" + DB)
	defer db.Close()

	if (err != nil) {
		log.Fatal("Error connecting database.")
		return false
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)

	if (err != nil) {
		log.Fatal("Error connecting database.")
	}

	fmt.Println("Mysql version is ", version)
	return true
}

func insertStudent(newStudent *student) {
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
	if err != nil {
		log.Fatal("Query executed with errors!")
	} else {
		for rows.Next() {
			// var s student
			var name string
			var id int
			var percentage float32
			var age byte
			rows.Scan(&id, &name, &age, &percentage)
			students = append(students, student{id, name, age, percentage})
		}
	}
	return students
}

func fetchStudents() []student {
	db, err := sql.Open(CLIENT, fmt.Sprintf("%s:%s@/%s", USER, PWD, DB))
	defer db.Close()

	if err != nil {
		log.Fatal("Error connecting database!")
	}

	fetchQuery := fmt.Sprintf("SELECT Id, Name, Age, Percentage FROM Student") 
	fmt.Println(fetchQuery)
	
	var students []student
	rows, err := db.Query(fetchQuery)
	if err != nil {
		log.Fatal("Query executed with errors!")
	} else {
		for rows.Next() {
			var name string
			var id int
			var percentage float32
			var age byte
			rows.Scan(&id, &name, &age, &percentage)
			students = append(students, student{id, name, age, percentage})
		}
	}
	return students
}

func updateStudent(updatedStudent *student) {
	db, err := sql.Open(CLIENT, fmt.Sprintf("%s:%s@/%s", USER, PWD, DB))
	if err != nil {
		log.Fatal("Error connecting database!")
		return
	}

	updateQuery := fmt.Sprintf("UPDATE Student SET name='%s', age=%d, percentage=%f WHERE id=%d", (*updatedStudent).name, (*updatedStudent).age, (*updatedStudent).percentage, (*updatedStudent).id)
	res, err := db.Exec(updateQuery)

	if err != nil {
		log.Fatal("Query executed with errors!")
		return
	}

	if rows, _ := res.RowsAffected(); rows == 1 {
		fmt.Println("Student updated.")
	}
}

func deleteStudent(id int) {
	db, err := sql.Open(CLIENT, fmt.Sprintf("%s:%s@/%s", USER, PWD, DB))
	if err != nil {
		log.Fatal("Error connecting database!")
		return
	}

	updateQuery := fmt.Sprintf("DELETE FROM Student WHERE id=%d", id)
	res, err := db.Exec(updateQuery)

	if err != nil {
		log.Fatal("Query executed with errors!")
		return
	}

	if rows, _ := res.RowsAffected(); rows == 1 {
		fmt.Println("Student deleted.")
	}
}

func showStudentsForUpdateAndDelete(students []student) int {
	fmt.Println("ID\tNAME\tAGE\tPERCENTAGE")
	for _, s := range(students) {
		fmt.Printf("%d\t%s\t%d\t%f\n", s.id, s.name, s.age, s.percentage)
	}

	var id int
	fmt.Print("Select any one student id displayed above: ")
	fmt.Scanln(&id)

	var fail bool = false
	for _, s := range(students) {
		if id != s.id {
			fail = true
		}
	}
	if fail {
		fmt.Println("Invalid id")
		return -1
	}

	return id
}


func main() {
	var choice int8
	for {
		if(!verifyConnecton()) {
			fmt.Println("Couldn't reach database server!")
			return
		}
		fmt.Println("Select any one of the following")
		fmt.Println("1. Add a student")
		fmt.Println("2. Display all students")
		fmt.Println("3. Search student")
		fmt.Println("4. Update student")
		fmt.Println("5. Delete student")
		fmt.Println("6. Exit")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var newStudent student
			fmt.Print("Enter student name, age and percentage respectively: ")
			fmt.Scan(&newStudent.name, &newStudent.age, &newStudent.percentage)
			insertStudent(&newStudent)
		
		case 2:
			students := fetchStudents()
			if len(students) > 0 {
				fmt.Println("ID\tNAME\tAGE\tPERCENTAGE")
				for _, s := range(students) {
					fmt.Printf("%d\t%s\t%d\t%.2f\n", s.id, s.name, s.age, s.percentage)
				}
			} else {
				fmt.Println("No data!")
			}
		case 3:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter student name: ")
			scanner.Scan()
			name := scanner.Text()
			students := fetchStudent(name)
			if len(students) > 0 {
				fmt.Println("ID\tNAME\tAGE\tPERCENTAGE")
				for _, s := range(students) {
					fmt.Printf("%d\t%s\t%d\t%f\n", s.id, s.name, s.age, s.percentage)
				}
			} else {
				fmt.Println("No data!")
			}
		case 4:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter student name: ")
			scanner.Scan()
			name := scanner.Text()
			students := fetchStudent(name)
			if len(students) > 0 {
				id := showStudentsForUpdateAndDelete(students)
				if id > 0 {
					var newStudent student
					fmt.Print("Enter student name, age and percentage respectively: ")
					fmt.Scanln(&newStudent.name, &newStudent.age, &newStudent.percentage)
					newStudent.id = id
					updateStudent(&newStudent)
				}
			} else {
				fmt.Println("No data!")
			}

		case 5:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter student name: ")
			scanner.Scan()
			name := scanner.Text()
			students := fetchStudent(name)
			if len(students) > 0 {
				id := showStudentsForUpdateAndDelete(students)
				if id > 0 {
					deleteStudent(id)
				}
			} else {
				fmt.Println("No data!")
			}
		default: return
		}
	}

}

