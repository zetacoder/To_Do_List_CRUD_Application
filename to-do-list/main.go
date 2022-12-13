package main

import (
	"bufio"
	"database/sql" // Contains funcs to open, prepare, query, manipulate the database and more.
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // A MySQL-Driver for Go's database/sql package
)

var selected_option int

var newTask Task

// checkErr checks if there is an error and panics if it is the case.
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Task struct {
	ID          int
	Name_task   string
	Description string
	Responsable string
	Completed   bool
}

// CREATE storage new inputs from the user and create new task.
func Create(t Task) (e error) {
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)

	defer db.Close()

	// Selected database.
	_, err = db.Exec("USE to_do_list")
	checkErr(err)

	// We prepare for the input in order to prevent SQL injections.
	prepareSentence, err := db.Prepare("INSERT INTO tasks (ID, name_task, description, responsable, completed) VALUES(?,?, ?,?,?)")
	checkErr(err)

	defer prepareSentence.Close()
	// Execute sentence for every '?'
	_, err = prepareSentence.Exec(t.ID, t.Name_task, t.Description, t.Responsable, t.Completed)
	checkErr(err)

	return nil
}

// READ selects rows from the MySQL table "to_do_list", reads it, scans it and returns a slice of the tasks.
func Read() ([]Task, error) {

	var tasks []Task

	// Open db.
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)
	defer db.Close()

	// Selected database.
	_, err = db.Exec("USE to_do_list")
	checkErr(err)

	// We execute a query iterating over all the rows.
	rows, err := db.Query("SELECT * FROM tasks")
	checkErr(err)
	defer rows.Close()

	var t Task

	// We read all the rows asigning the values to t (type Task))
	// Then, append the values to tasks and returning it.
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Name_task, &t.Description, &t.Responsable, &t.Completed)
		checkErr(err)
		tasks = append(tasks, t)
	}

	return tasks, nil

}

// UPDATE selects some task and modifies it.
func Update(n Task) error {

	// Connected to MySQL database.
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/to_do_list")
	checkErr(err)

	// Close db after use.
	defer db.Close()

	preparedSentence, err := db.Prepare("UPDATE tasks SET name_task = ?, description = ?, responsable = ?, completed = ? WHERE id = ?")
	checkErr(err)

	defer preparedSentence.Close()

	_, err = preparedSentence.Exec(n.Name_task, n.Description, n.Responsable, n.Completed, n.ID)
	checkErr(err)
	return nil
}

// DELETE selects some task and remove it from the db.
func Delete(n Task) error {
	// Connected to MySQL database.
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/to_do_list")
	checkErr(err)

	// Close db after use.
	defer db.Close()

	preparedSentence, err := db.Prepare("DELETE FROM tasks WHERE ID = ?")
	checkErr(err)

	defer preparedSentence.Close()

	_, err = preparedSentence.Exec(n.ID)
	checkErr(err)

	return nil
}

func main() {

	// Connected to MySQL database.
	db, err := sql.Open("mysql", "root:<yourpassword>.@tcp(127.0.0.1:3306)/")
	checkErr(err)

	// Close db after use.
	defer db.Close()

	// Created To Do List Database.
	_, err = db.Exec("CREATE DATABASE if not exists to_do_list")
	checkErr(err)

	// Selected database.
	_, err = db.Exec("USE to_do_list")
	checkErr(err)

	// Created table where the data is going to storaged and retrieved.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (ID tinyint NOT NULL auto_increment, name_task tinytext NOT NULL, description tinytext NOT NULL, responsable tinytext NOT NULL, completed bool not null, primary key (ID))")
	checkErr(err)

	fmt.Println("")
	fmt.Println("WELCOME TO TO_DO_LIST APP BY ZETACODER.")
	fmt.Println("CREATE NEW TASKS TO DO, MODIFY IT AND DELETE IT ONCE COMPLETED.ENJOY IT :D")
	menu := `
What do you want to do: 
(1) NEW TASK.
(2) SHOW ALL.
(3) MODIFY.
(4) COMPLETE.
(5) DELETE.
(6) EXIT APP.
`

	// Show menu and depending of the option, execute the CRUD funcs.
	for selected_option != 6 {
		fmt.Println("")
		fmt.Println(menu)
		fmt.Scan(&selected_option)
		scanner := bufio.NewScanner(os.Stdin)

		switch selected_option {
		// Create new task
		case 1:
			fmt.Println("Name of the task:")
			fmt.Scanln()
			if scanner.Scan() {
				newTask.Name_task = scanner.Text()
			}
			fmt.Println("")

			fmt.Println("Description:")
			if scanner.Scan() {
				newTask.Description = scanner.Text()
			}

			fmt.Println("Person responsable:")
			time.Sleep(1 * time.Second)
			if scanner.Scan() {
				newTask.Responsable = scanner.Text()
			}

			Create(newTask)
			time.Sleep(1 * time.Second)
			fmt.Println("TASK ADDED SUCCESFULLY!")
			time.Sleep(1 * time.Second)
		// Show all tasks
		case 2:
			tasks, err := Read()
			if err != nil {
				fmt.Printf("Cant get tasks: %v", err)
			} else {
				for _, task := range tasks {
					fmt.Println("====================")
					fmt.Printf("ID: %d\n", task.ID)
					fmt.Printf("Name of task: %s\n", task.Name_task)
					fmt.Printf("Description: %s\n", task.Description)
					fmt.Printf("Responsable: %s\n", task.Responsable)
					fmt.Printf("Completed: %v\n", task.Completed)
				}
			}
		case 3:
			fmt.Println("Choose the ID of the task to modify:")
			fmt.Scanln()
			fmt.Scanln(&newTask.ID)
			fmt.Println("Enter new task name:")
			if scanner.Scan() {
				newTask.Name_task = scanner.Text()
			}
			fmt.Println("Enter new description:")
			if scanner.Scan() {
				newTask.Description = scanner.Text()
			}
			fmt.Println("Enter the new responsable of the task:")
			if scanner.Scan() {
				newTask.Responsable = scanner.Text()
			}
			err := Update(newTask)
			if err != nil {
				fmt.Printf("Error updating: %v", err)
			} else {
				time.Sleep(1 * time.Second)
				fmt.Println("TASK UPDATED SUCCESFULLY")
				time.Sleep(1 * time.Second)
			}
		case 4:
			fmt.Println("Choose the ID of the task you completed:")
			fmt.Scanln()
			fmt.Scanln(&newTask.ID)
			Delete(newTask) // Mark as completed

			fmt.Printf("Congratulations! Task %d marked as completed", newTask.ID)
			time.Sleep(1 * time.Second)
		case 5:
			fmt.Println("Choose the ID of the task you want to remove:")
			fmt.Scanln()
			fmt.Scanln(&newTask.ID)
			Delete(newTask) // task deleted

			fmt.Println("Task removed from the list.")
			time.Sleep(1 * time.Second)
		case 6:
			fmt.Println("See ya!")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}
}
