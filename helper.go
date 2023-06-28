package main

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

const ( // color
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
)

func Red(msg string) string {
	return colorRed + msg + colorReset
}
func Green(msg string) string {
	return colorGreen + msg + colorReset
}
func Yellow(msg string) string {
	return colorYellow + msg + colorReset
}

// chatComplete send request and get response from the OpenAI
// it uses 'gpt-3.5-turbo'
func chatCompleteStream(messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	stream, err := openaiClient.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     GPTModel,
			Messages:  messages,
			MaxTokens: 300,
			Stream:    true,
		},
	)
	return stream, err
}

// command
// func helpMessage() string {
// 	help := colorStr(Green, "help")
// 	config := colorStr(Green, "config")
// 	context := colorStr(Green, "context")
// 	reset := colorStr(Green, "reset")
// 	clear := colorStr(Green, "clear")
// 	exit := colorStr(Green, "exit")
// 	q := colorStr(Green, "q")

// 	return fmt.Sprintf(`Usage:
//   - %s - Displays this help message.
//   - %s - Displays configuration information.
//   - %s - Displays the conversation context which reserved at the moment.
//   - %s - Reset all the conversation context.
//   - %s - Clear terminal.
//   - %s or %s - Exits the app.
// 	`, help, config, context, reset, clear, exit, q)
// }
