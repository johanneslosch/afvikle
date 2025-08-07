package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func main() {
	cli := clir.NewCli("afv", "Short for afvikle. CLI to speed up the process of running multiple scripts without creating another script. Run from anywhere.", "v1.0.0")

	var commands string

	cli.NewSubCommand("list", "Returns a list of commands runnable with afvikle").
		Action(func() error {
			fmt.Println(commands)
			return nil
		})

	// Starte die CLI
	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}