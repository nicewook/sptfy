package main

import (
	"context"

	"github.com/sashabaranov/go-openai"
)



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
