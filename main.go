package main

import (
        "fmt"

        "github.com/manifoldco/promptui"
)

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

        fmt.Printf("You choose %q\n", res)
}
