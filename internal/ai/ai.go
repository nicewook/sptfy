package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nicewook/sptfy/internal/config"
	"github.com/nicewook/sptfy/internal/sp"
	"github.com/sashabaranov/go-openai"
)

const (
	GPTModel = "gpt-3.5-turbo-0613"
)

var openaiClient *openai.Client

func init() {
	config.InitConfig()
	openaiClient = openai.NewClient(config.GetConfig().OpenAIAPIKey)
}

func GeneratePlaylist(prompt string, num int) (funcName string, pl sp.Playlist) {

	fmt.Printf("generating play list of %d tracks\n", num)
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a helpful playlist generating assistant. 
					      You should generate a list of songs and their artists according to a text prompt.`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Generate a playlist of %d songs based on the prompt: ###%s### ", num, prompt),
		},
	}

	resp, err := chatComplete(messages)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Choices[0].FinishReason != openai.FinishReasonFunctionCall {
		log.Println("functions call is not activatied.")
	}
	fmt.Println("playlist generated.")

	funcName = resp.Choices[0].Message.FunctionCall.Name
	funcArg := resp.Choices[0].Message.FunctionCall.Arguments

	if err := json.Unmarshal([]byte(funcArg), &pl); err != nil {
		log.Println(err)
		return "", pl
	}
	return funcName, pl
}

var msg = json.RawMessage(`
	{
		"type": "object",
		"properties": {
			"playlist": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"song": {
							"type": "string",
							"description": "song title"
						},
						"artist": {
							"type": "string",
							"description": "artist or group name"
						}
					},
					"required": ["song", "artist"]
				}
			}
		}
	}
`)

func chatComplete(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     GPTModel,
			Messages:  messages,
			MaxTokens: 300,
			Functions: []openai.FunctionDefinition{
				{
					Name:       "SpotifyPlaylistGenerator",
					Parameters: &msg,
				},
			},
			FunctionCall: "auto",
		},
	)
	return resp, err
}

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
