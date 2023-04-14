package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Title       string
	Description string
	DueDate     time.Time
	Priority    int
	Status      bool
}

type TaskList struct {
	Tasks []Task
}

func (tl *TaskList) addTask(task Task) {
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) removeTask(title string) {
	for i, task := range tl.Tasks {
		if task.Title == title {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return
		}
	}
}

func (tl *TaskList) showTask(title string) {
	for _, task := range tl.Tasks {
		if task.Title == title {
			fmt.Println("Title:", task.Title)
			fmt.Println("Description:", task.Description)
			fmt.Println("Due Date:", task.DueDate.Format("02 January"))
			fmt.Println("Priority:", task.Priority)
			fmt.Println("Status:", task.Status)
		}
	}
}

func (tl *TaskList) displayTasks() {
	fmt.Println("==========================================================")
	fmt.Println("Title         || Due Date        || Priority || Status")
	fmt.Println("----------------------------------------------------------")
	for _, task := range tl.Tasks {
		fmt.Printf("%-14s || %-15s || %-8d || %t\n", task.Title, task.DueDate.Format("02 January"), task.Priority, task.Status)
	}
	fmt.Println("==========================================================")
	fmt.Println()
}

func (tl *TaskList) sortByDueDate() {
	sort.SliceStable(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].DueDate.Before(tl.Tasks[j].DueDate)
	})
}

func (tl *TaskList) sortByPriority() {
	sort.SliceStable(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].Priority < tl.Tasks[j].Priority
	})
}

func initializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		panic(err)
	}

	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS tasks (
		title TEXT PRIMARY KEY,
		description TEXT,
		duedate DATETIME,
		priority INTEGER,
		status BOOLEAN
	);`)
	statement.Exec()

	return db
}

func saveTasksToDatabase(db *sql.DB, tasks []Task) {
	statement, _ := db.Prepare("DELETE FROM tasks")
	statement.Exec()

	for _, task := range tasks {
		statement, _ := db.Prepare("INSERT INTO tasks (title, description, duedate, priority, status) VALUES (?, ?, ?, ?, ?)")
		statement.Exec(task.Title, task.Description, task.DueDate, task.Priority, task.Status)
	}
}

func removeTaskFromDatabase(db *sql.DB, title string) error {
	statement, err := db.Prepare("DELETE FROM tasks WHERE title = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(title)
	if err != nil {
		return err
	}

	return nil
}

func loadTasksFromDatabase(db *sql.DB) []Task {
	rows, _ := db.Query("SELECT title, description, duedate, priority, status FROM tasks")
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var title, description string
		var dueDate time.Time
		var priority int
		var status bool

		rows.Scan(&title, &description, &dueDate, &priority, &status)
		task := Task{Title: title, Description: description, DueDate: dueDate, Priority: priority, Status: status}
		tasks = append(tasks, task)
	}

	return tasks
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	db := initializeDatabase()
	defer db.Close()

	taskList := TaskList{Tasks: loadTasksFromDatabase(db)}

	fmt.Println("Welcome to the TODO list manager!")
	fmt.Println("Please enter a command (add, complete, remove, show, sortDate, sortPriority or exit):")

	for {
		taskList.displayTasks()
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		switch command {
		case "add":
			fmt.Print("Enter task title: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Enter task description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)

			fmt.Print("Enter task due date (YYYY-MM-DD): ")
			dateString, _ := reader.ReadString('\n')
			dateString = strings.TrimSpace(dateString)
			dueDate, _ := time.Parse("2006-01-02", dateString)

			fmt.Print("Enter task priority: ")
			var priority int
			fmt.Scanf("%d\n", &priority)

			task := Task{
				Title:       title,
				Description: description,
				DueDate:     dueDate,
				Priority:    priority,
				Status:      false,
			}

			taskList.addTask(task)
			saveTasksToDatabase(db, taskList.Tasks)
			fmt.Println("Task added successfully!")

		case "complete":
			fmt.Print("Enter the title of the task you want to mark as completed: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			for i, task := range taskList.Tasks {
				if task.Title == title {
					taskList.Tasks[i].Status = true
					fmt.Println("Task marked as completed!")
					break
				}
			}

		case "remove":
			fmt.Print("Enter the title of the task you want to remove: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			err := removeTaskFromDatabase(db, title)
			if err != nil {
				fmt.Println("Error removing task:", err)
			} else {
				taskList.removeTask(title)
				fmt.Println("Task removed successfully!")
			}

		case "show":

			fmt.Print("Enter the title of the task you want to be shown: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			taskList.showTask(title)

		case "sortDate":
			taskList.sortByDueDate()

		case "sortPriority":
			taskList.sortByPriority()

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid command. Please enter a valid command (add, complete, remove, show, or exit).")
		}
	}
}
