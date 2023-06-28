package main

import (
	"log"

	"github.com/chzyer/readline"
)

func getReadline() *readline.Instance {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          Green(">> "),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	return rl
}
