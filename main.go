package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
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

func (tl *TaskList) showTask() {
	for _, task := range tl.Tasks {
		fmt.Println("Title:", task.Title)
		fmt.Println("Description:", task.Description)
		fmt.Println("Due Date:", task.DueDate)
		fmt.Println("Priority:", task.Priority)
		fmt.Println("Status:", task.Status)
		fmt.Println()
	}
}

func (tl *TaskList) displayTasks() {
	fmt.Println("Title         || Due Date        || Priority || Status")
	fmt.Println("--------------------------------------------------")
	for _, task := range tl.Tasks {
		fmt.Printf("%-14s || %-15s || %-8d || %t\n", task.Title, task.DueDate.Format("02 January"), task.Priority, task.Status)
	}
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	taskList := TaskList{}

	task1 := Task{
		Title:       "Task 1",
		Description: "This is task 1",
		DueDate:     time.Now().AddDate(0, 0, 7),
		Priority:    2,
		Status:      false,
	}

	task2 := Task{
		Title:       "Task 2",
		Description: "This is task 2",
		DueDate:     time.Now().AddDate(0, 0, 3),
		Priority:    1,
		Status:      false,
	}

	taskList.addTask(task1)
	taskList.addTask(task2)

	// taskList.showTasks()

	// taskList.sortByDueDate()
	// fmt.Println("Sorted by due date:")
	// taskList.showTasks()

	// taskList.sortByPriority()
	// fmt.Println("Sorted by priority:")
	// taskList.showTasks()

	// taskList.removeTask("Task 1")
	// fmt.Println("Task 1 removed:")
	// taskList.showTasks()

	fmt.Println("Welcome to the TODO list manager!")
	fmt.Println("Please enter a command (add, complete, remove, show, or exit):")

	for {
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

			taskList.removeTask(title)
			fmt.Println("Task removed successfully!")

		case "show":
			taskList.displayTasks()

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid command. Please enter a valid command (add, complete, remove, show, or exit).")
		}
	}
}
