package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

type Task struct {
	Id          int       `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("no args")
	}
	cmd := args[1]
	switch cmd {
	case "add":
		if len(args) < 3 {
			log.Fatalln("usage: add \"task detail\"")
		}
		content := args[2]
		addTask(content)
	case "update":
		if len(args) < 4 {
			log.Fatalln("usage: update id \"update detail\"")
		}
		id := args[2]
		content := args[3]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln("can't parse id")
		}
		updateTask(taskId, content)
	case "delete":
		if len(args) < 3 {
			log.Fatalln("usage: delete task_id")
		}
		id := args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln("can't parse id")
		}
		deleteTask(taskId)
	case "list":
		listTask(args[2:]...)
	case "mark-done":
		if len(args) < 3 {
			log.Fatalln("usage: mark-done <task_id>")
		}
		id := args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln("can't parse id")
		}
		markTask(taskId, "done")
	case "mark-in-progress":
		if len(args) < 3 {
			log.Fatalln("usage: mark-in-progress <task_id>")
		}
		id := args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalln("can't parse id")
		}
		markTask(taskId, "in-progress")
	default:
		log.Fatalln("no match command. available are: add, update, delete, list, mark-done. mark-in-progress")
	}
}

func markTask(id int, status string) {
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalln(err)
	}
	found := false
	for i, task := range tasks {
		if task.Id == id {
			found = true
			task.Status = status
			task.UpdatedAt = time.Now()
			tasks[i] = task
			break
		}
	}

	if !found {
		log.Fatalf("Task with ID %d not found\n", id)
	} else {
		saveTask(&tasks)
		log.Printf("Task marked as %s.", status)
	}
}

func timeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff.Seconds() < 60 {
		return "few minute ago"
	} else if diff.Minutes() < 60 {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff.Hours() < 24 {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}
}

func listTask(statuses ...string) {
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalln(err)
	}
	status := "all"
	if len(statuses) > 0 {
		status = statuses[0]
	}
	var filteredTasks []Task
	if status == "all" {
		filteredTasks = tasks
	} else {
		newTasks := slices.DeleteFunc(tasks, func(task Task) bool {
			return task.Status != status
		})
		filteredTasks = newTasks
	}
	if len(filteredTasks) <= 0 {
		log.Fatalln("No tasks found")
	} else {
		for _, task := range filteredTasks {
			fmt.Printf("[%d] %s (%s) - %s\n", task.Id, task.Description, task.Status, timeAgo(task.UpdatedAt))
		}
	}
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalln(err)
	}
	index := slices.IndexFunc(tasks, func(task Task) bool {
		return task.Id == id
	})

	if index != -1 {
		tasks = slices.Delete(tasks, index, index+1)
		saveTask(&tasks)
		log.Println("Task deleted successfully")
	} else {
		log.Printf("Task with ID %d not found", id)
	}
}

func updateTask(id int, content string) {
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalln(err)
	}
	found := false
	for i, task := range tasks {
		if task.Id == id {
			found = true
			task.Description = content
			task.UpdatedAt = time.Now()
			tasks[i] = task
			break
		}
	}

	if !found {
		log.Fatalf("Task with ID %d not found\n", id)
	} else {
		saveTask(&tasks)
		log.Println("Task updated successfully.")
	}
}

func addTask(content string) {
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalln(err)
	}
	var id int
	if len(tasks) > 0 {
		lastTask := tasks[len(tasks)-1]
		id = lastTask.Id + 1
	} else {
		id = 1
	}

	task := Task{
		Id:          id,
		Description: content,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, task)

	saveTask(&tasks)
	log.Printf("Task added successfully (ID: %d)\n", id)
}

func loadTasks() ([]Task, error) {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil, errors.New("can't read tasks.json file")
	}
	var tasks []Task
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return []Task{}, nil
	}
	return tasks, nil
}

func saveTask(tasks *[]Task) {
	bytes, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatalln("can't save file")
	}
	err = os.WriteFile("tasks.json", bytes, os.ModePerm)
	if err != nil {
		log.Fatalln("can't write file")
	}
}
