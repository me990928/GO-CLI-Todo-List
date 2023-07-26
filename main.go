package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/manifoldco/promptui"
)

type Todo struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	StartDate string `json:"stDate"`
	EndDate   string `json:"edDate"`
}

func main() {
	prompt := promptui.Select{
		Label: "Menu",
		Items: []string{"Add Todo", "Show Todo"},
	}

	_, res, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if res == "Add Todo" {
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
		fmt.Scan(&date)
		var endDate = time.Now().AddDate(0, 0, date).Format("2006-01-02")
		fmt.Printf("endDate:%s", endDate)

		prompt := promptui.Select{
			Label: "What to save?",
			Items: []string{"Yes", "No"},
		}

		_, res, err := prompt.Run()

		if err != nil {
			fmt.Println("err")
		}

		if res == "Yes" {
			todo := []Todo{
				{Title: title, Body: body, StartDate: startDate, EndDate: endDate},
			}

	                fmt.Println("Writing Start!")
			writeJson(todo)

		} else {
			main()
		}
	}
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

	fmt.Println("Write Todo Success!")

}
