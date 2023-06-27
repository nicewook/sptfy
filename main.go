package main

import (
	"fmt"
	"io"
	"log"

	"github.com/chzyer/readline"
	"github.com/nicewook/sptfy/internal/config"
	"github.com/sashabaranov/go-openai"
)

func getReadline() *readline.Instance {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}

	return rl
}

func main() {
	log.Println("sptfy")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg

	client = openai.NewClient(cfg.OpenAIAPIKey)
	rl := getReadline()
	defer rl.Close()
	rl.CaptureExitSignal()

	for {
		prompt, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(prompt) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		generatePlaylist(prompt, 8)
	}

}

func generatePlaylist(prompt string, num int) {
	log.Println("qurate play list")
	exampleResponse := `
	{
		"songs": [
			{"song": "Everybody Hurts", "artist": "R.E.M."},
			{"song": "Nothing Compares 2 U", "artist": "Sinead O'Connor"},
			{"song": "Tears in Heaven", "artist": "Eric Clapton"},
			{"song": "Hurt", "artist": "Johnny Cash"},
			{"song": "Yesterday", "artist": "The Beatles"}
		]
	}
	`
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem, Content: `You are a helpful playlist generating assistant. 
					You should generate a list of songs and their artists according to a text prompt.
					`,
		},
		{Role: openai.ChatMessageRoleUser, Content: "Generate a playlist of 5 songs based on this prompt: super super sad songs"},
		{Role: openai.ChatMessageRoleAssistant, Content: exampleResponse},
		{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Generate a playlist of %d songs based on this prompt: %s", num, prompt)},
	}

	resp, err := chatComplete(messages)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("got it")
	log.Println("len choice:", len(resp.Choices))
	fmt.Printf("%+v", resp.Choices[0].Message.Content)
	fmt.Printf("%+v", resp.Choices[0].Message.FunctionCall.Name)
	fmt.Printf("%+v", resp.Choices[0].Message.FunctionCall.Arguments)
	// fmt.Printf("%+v", resp.Choices[0].Content.FunctionCall)
	// if len(resp.Functions) > 0 {
	// 	var fcb []byte
	// 	b := resp.Functions[0].Parameters
	// 	fcb, err = json.Marshal(b)
	// 	if err != nil {
	// 		log.Println("err:", err)
	// 		return
	// 	}
	// }
	// fmt.Println(string(fcb))

}
