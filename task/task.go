package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

const dataFile = "data.json"

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Task struct {
	ID          int       `json:"id"`          //Уникальный идентификатор задачи
	Description string    `json:"description"` // Описание задачи
	Status      string    `json:"status"`      // Статус задачи (todo, in-progress, done)
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func load() ([]*Task, error) {
	var tasks []*Task
	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, fmt.Errorf("ошибка чтения файла %s: %w", dataFile, err)
	}
	if len(file) == 0 {
		return tasks, nil
	}
	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, fmt.Errorf("не удалось разобрать файл %s: %w", dataFile, err)
	}
	return tasks, nil
}

func save(tasks []*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	tmp := dataFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmp, dataFile); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func NewTask(title string) error {
	if title == "" {
		return errors.New("описание задачи пустое")
	}
	tasks, err := load()
	if err != nil {
		return err
	}
	var maxID int
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1
	task := &Task{
		ID:          newID,
		Description: title,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, task)
	return save(tasks)
}

func UpdateTask(id int, title string) error {
	if title == "" {
		return errors.New("описание задачи пустое")
	}
	tasks, err := load()
	if err != nil {
		return err
	}
	idTask := findTaskByID(id, tasks)
	if idTask == -1 {
		return errors.New("неверно указан id")
	}
	tasks[idTask].Description = title
	tasks[idTask].UpdatedAt = time.Now()
	return save(tasks)
}

func findTaskByID(id int, tasks []*Task) int {
	for i, task := range tasks {
		if task == nil {
			continue
		}
		if task.ID == id {
			return i
		}
	}
	return -1
}

func DeleteTask(id int) error {
	tasks, err := load()
	if err != nil {
		return err
	}
	idTask := findTaskByID(id, tasks)
	if idTask == -1 {
		return errors.New("задача не найдена")
	}
	tasks = append(tasks[:idTask], tasks[idTask+1:]...)
	return save(tasks)
}

func MarkInProgress(id int) error {
	tasks, err := load()
	if err != nil {
		return err
	}
	idTask := findTaskByID(id, tasks)
	if idTask == -1 {
		return errors.New("задача не найдена")
	}
	tasks[idTask].Status = StatusInProgress
	tasks[idTask].UpdatedAt = time.Now()
	return save(tasks)
}

func MarkDone(id int) error {
	tasks, err := load()
	if err != nil {
		return err
	}
	idTask := findTaskByID(id, tasks)
	if idTask == -1 {
		return errors.New("задача не найдена")
	}
	tasks[idTask].Status = StatusDone
	tasks[idTask].UpdatedAt = time.Now()
	return save(tasks)
}

func ListByStatus(status string) {
	tasks, err := load()
	if err != nil {
		color.Red("Ошибка при загрузке задач: %v", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("Задач нет")
		return
	}
	found := false
	for _, task := range tasks {
		if task == nil {
			continue
		}
		if status != "" && task.Status != status {
			continue
		}
		found = true
		fmt.Println()
		label1 := color.MagentaString("Id задачи: ")
		label2 := color.WhiteString("%d", task.ID)
		fmt.Println(label1 + label2)
		label1 = color.CyanString("Описание: ")
		label2 = color.WhiteString("%s", task.Description)
		fmt.Println(label1 + label2)
		label1 = color.CyanString("Статус задачи: ")
		label2 = color.WhiteString("%s", task.Status)
		fmt.Println(label1 + label2)
		label1 = color.CyanString("Время создания: ")
		label2 = color.WhiteString("%s", task.CreatedAt.Format(time.RFC3339))
		fmt.Println(label1 + label2)
		label1 = color.CyanString("Время обновления: ")
		label2 = color.WhiteString("%s", task.UpdatedAt.Format(time.RFC3339))
		fmt.Println(label1 + label2)
	}
	if !found {
		if status == "" {
			fmt.Println("Задач нет")
		} else {
			fmt.Printf("Задач со статусом \"%s\" не найдено\n", status)
		}
	}
}

func List() {
	ListByStatus("")
}

func ListDone() {
	ListByStatus(StatusDone)
}

func ListTodo() {
	ListByStatus(StatusTodo)
}

func ListInProgress() {
	ListByStatus(StatusInProgress)
}
