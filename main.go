package main

import (
	"TaskTracker/task"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	color.White("____Трекер задач____")
	if len(os.Args) < 2 {
		color.Red("Команда не указана")
		printCommand()
		return
	}
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			color.Red("Неверное количество аргументов")
			return
		}
		title := strings.Join(os.Args[2:], " ")
		if err := task.NewTask(title); err != nil {
			color.Red("Ошибка: %v", err)
		} else {
			color.Green("Задача добавлена")
		}
	case "list":
		switch len(os.Args) {
		case 2:
			task.List()
		case 3:
			arg := os.Args[2]
			switch arg {
			case "todo":
				fmt.Println("Задачи которые нужно выполнить")
				task.ListTodo()
			case "done":
				fmt.Println("Выполненныe задачи")
				task.ListDone()
			case "in-progress":
				fmt.Println("Задачи находящиecя в процессе выполнения")
				task.ListInProgress()
			default:
				color.Red("Неверно введенная команда")
				printCommand()
			}
		default:
			color.Red("Слишком много аргументов для list")
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("update id \"описание\" - Обновить задачу")
			return
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			color.Red("Ошибка: %v", err)
			return
		}
		title := strings.Join(os.Args[3:], " ")
		err = task.UpdateTask(id, title)
		if err != nil {
			color.Red("Ошибка: %v", err)
		} else {
			color.Green("Задача c id %d обновлена", id)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("delete id - Удалить задачу")
			return
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			color.Red("Ошибка: %v", err)
			return
		}
		err = task.DeleteTask(id)
		if err != nil {
			color.Red("Ошибка: %v", err)
		} else {
			color.Green("Задача удалена")
		}
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("mark-in-progress id - Сменить статус задачи на статус \"в работе\"")
			return
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			color.Red("Ошибка: %v", err)
			return
		}
		if err := task.MarkInProgress(id); err != nil {
			color.Red("Ошибка: %v", err)
		} else {
			color.Green("Задача %d в работе", id)
		}
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("mark-done id - Сменить статус задачи на статус \"выполнена\"")
			return
		}
		id, err := parseID(os.Args[2])
		if err != nil {
			color.Red("Ошибка: %v", err)
			return
		}
		if err := task.MarkDone(id); err != nil {
			color.Red("Ошибка: %v", err)
		} else {
			color.Green("Задача %d выполнена", id)
		}
	default:
		printCommand()
	}
}

func printCommand() {
	fmt.Println(`
Команды:
list			- Список задач
add "описание"          - Добавить задачу
update id "описание"  	- Обновить задачу 
delete id              	- Удалить задачу 
mark-in-progress id   	- Сменить статус задачи на статус "в работе" 
mark-done id          	- Сменить статус задачи на статус "выполнена"
list todo            	- Задачи которые нужно выполнить
list done            	- Выполненные задачи
list in-progress     	- Задачи находящиеся в процессе выполнения`)
}

func parseID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return 0, fmt.Errorf("некорректный id: %w", err)
	}
	return id, nil
}
