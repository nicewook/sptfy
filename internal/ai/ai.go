package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nicewook/sptfy/internal/color"
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

// GeneratePlayListName generate playlist name from prompt and playlist
func GeneratePlayListName(prompt, playlist string) string {
	fmt.Println(color.Blue("Generating playlist name!"))

	loading := spinner.New([]string{".", "..", "...", "....", "....."}, 150*time.Millisecond)
	loading.Prefix = color.Yellow("loading")
	loading.Color("yellow")

	// fmt.Println(playlist)
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a helpful playlist naming assistant. 
					      You name accurate and artistic Spotify playlist name from th prompt which inspired generating playlist and the playlist itself.
								The name should not be over 10 words`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(`Generate a Spotify playlist name based on the prompt: ###sunny###, and generated playlist below
			{
				"playlist": [
					{
						"song": "Soul Food",
						"artist": "Goodie Mob"
					},
					{
						"song": "My Music",
						"artist": "Soulja Boy"
					},
					{
						"song": "Food For My Soul",
						"artist": "Jhene Aiko"
					},
					{
						"song": "Music Saved My Life",
						"artist": "Joey Bada$$"
					}
				]
			}

      Output format is only generated playlist name itself`, prompt, playlist),
		},
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: "Sunny day afternoon, with music",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(`Generate a Spotify playlist name based on the prompt: ###%s###, and generated playlist below\n%s
			Output format is only generated playlist name itself`, prompt, playlist),
		},
	}

	loading.Start()
	defer func() {
		if loading.Active() {
			loading.Stop()
		}
	}()
	resp, err := chatComplete(messages, false)
	if err != nil {
		log.Println(err)
		return prompt
	}
	loading.Stop()
	if resp.Choices[0].FinishReason != openai.FinishReasonStop {
		return prompt
	}
	
	playlistName := resp.Choices[0].Message.Content
	return playlistName

}

func GeneratePlaylist(prompt string, num int) (funcName string, pl sp.Playlist) {

	fmt.Println(color.Blue(fmt.Sprintf("Generating playlist of %d(or less) tracks!", num)))
	loading := spinner.New([]string{".", "..", "...", "....", "....."}, 150*time.Millisecond)
	loading.Prefix = color.Yellow("loading")
	loading.Color("yellow")

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

	loading.Start()
	defer func() {
		if loading.Active() {
			loading.Stop()
		}
	}()
	resp, err := chatComplete(messages, true)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Choices[0].FinishReason != openai.FinishReasonFunctionCall {
		log.Println("functions call is not activatied.")
		return
	}
	loading.Stop()
	fmt.Println("Playlist generated.")

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

func chatComplete(messages []openai.ChatCompletionMessage, useFunction bool) (openai.ChatCompletionResponse, error) {

	req := openai.ChatCompletionRequest{
		Model:     GPTModel,
		Messages:  messages,
		MaxTokens: 300,
	}

	if useFunction {
		req.Functions = []openai.FunctionDefinition{
			{
				Name:       "SpotifyPlaylistGenerator",
				Parameters: &msg,
			},
		}
		req.FunctionCall = "auto"
	}

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		req,
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
