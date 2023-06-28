package main

import (
	"log"

	"github.com/chzyer/readline"
	"github.com/nicewook/sptfy/internal/color"
)

func getReadline() *readline.Instance {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          color.Green(">> "),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	return rl
}
