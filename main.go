package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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
	default:
		log.Fatalln("no match command. available are: add")
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
