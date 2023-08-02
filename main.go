package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"time"

	"github.com/manifoldco/promptui"
	// "golang.org/x/text/search"
)

type Todo struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	StartDate string `json:"stDate"`
	EndDate   string `json:"edDate"`
	IsTodo    bool   `json:"isTodo"`
}

func main() {
	for {
		prompt := promptui.Select{
			Label: "Menu",
			Items: []string{"Add Todo", "Show Todo", "Exit"},
		}

		_, res, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			continue
		}

		switch res {
		case "Add Todo":
			addTodo()
		case "Show Todo":
			showTodo()
		case "Exit":
			fmt.Println("Bye...")
			return
		default:
			fmt.Println("Err...")
			return
		}

	}
}

func addTodo() {

	fmt.Printf("Title：")
	var title string
	fmt.Scan(&title)
	fmt.Printf("Body：")
	var body string
	fmt.Scan(&body)
	var startDate = time.Now().Format("2006-01-02")
	fmt.Printf("Start day： %s\n", startDate) // debug
	var date int = 0
	fmt.Printf("How many days?：")
	_, err := fmt.Scan(&date)
	if err != nil {
		fmt.Println("Non-numeric input confirmed.")
		return
	}
	var endDate = time.Now().AddDate(0, 0, date).Format("2006-01-02")
	fmt.Printf("endDate:%s", endDate)
	var isTodo = false

	template := &promptui.SelectTemplates{
		Label: "{{.}}",
		Details: `
Title: ` + title + `
Body: ` + body + `
Start day: ` + startDate + `
End Date: ` + endDate + `
Date: ` + fmt.Sprintf("%d", date),
	}

	prompt := promptui.Select{
		Label:     "What to save?",
		Templates: template,
		Items:     []string{"Yes", "No"},
	}

	_, res, err := prompt.Run()

	if err != nil {
		fmt.Println("err")
	}

	if res == "Yes" {
		todo := []Todo{
			{Title: title, Body: body, StartDate: startDate, EndDate: endDate, IsTodo: isTodo},
		}

		fmt.Println("Writing Start!")
		writeJson(todo)
		fmt.Println("Write Todo Success!")

	}

}

func showTodo() {
	var todos = []Todo{}
	todos, err := readJson()
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	dispAll(todos)
}

func writeJson(data []Todo) {

	beforeData := []Todo{}

	_, err := os.Stat("todo.json")

	if err != nil {
		var flag bool = true
		if os.IsNotExist(err) {
			// ファイルがない場合は新規作成
			// fmt.Println("Not found File")
			file, err := os.Create("todo.json")
			if err != nil {
				fmt.Println("Create json Error.")
				return
			}
			defer file.Close()
			flag = false
		} else {
			fmt.Println("Err：", err)
		}

		if flag {
			return
		}

	} else {

		file, err := os.Open("todo.json")

		if err != nil {
			fmt.Println("Open file error:", err)
			return
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(&beforeData)

		if err != nil {
			fmt.Println("err:", err)
			return
		}

		data = append(data, beforeData...)

	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Encoding Errors.")
		return
	}

	file, err := os.Create("todo.json")

	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

	_, err = file.Write(jsonData)

	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

}

func readJson() ([]Todo, error) {

	var todoData = []Todo{}

	file, err := os.Open("todo.json")
	if err != nil {
		return todoData, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return todoData, err
	}

	err = json.Unmarshal(data, &todoData)
	if err != nil {
		return todoData, err
	}

	return todoData, nil
}

func dispAll(todos []Todo) {

	const defDate = "2006-01-02"
	// todos[0].
	template := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "✔️ {{ .Title | green }}",
		Inactive: "{{ .Title}}",
		Selected: "✔️ {{ .Title | cyan }}",
		Details: `
============================
Body: {{.Body}}
StartDate: {{.StartDate}}
EndDate: {{.EndDate}}
IsTodo: {{.IsTodo}}`,
	}

	searcher := func(input string, index int) bool {
		todo := todos[index]
		name := strings.Replace(strings.ToLower(todo.Body), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Todo List",
		Items:     todos,
		Templates: template,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Err: ", err)
		return
	}

	prompt2 := promptui.Select{
		Label:     "Complete todo?",
		Templates: template,
		Items:     []string{"Yes", "No"},
	}

	_, res, err1 := prompt2.Run()

	if err1 != nil {
		fmt.Printf("Err: ", err1)
		return
	}

	if res == "Yes" {
		todos[i].IsTodo = true
		writeJson(todos)
	} else {
		todos[i].IsTodo = false
		writeJson(todos)
	}

	fmt.Println(res)

	return
}
