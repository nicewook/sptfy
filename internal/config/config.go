package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/chzyer/readline"
	"github.com/nicewook/sptfy/internal/color"
	"github.com/spf13/viper"
)

type Config struct {
	OpenAIAPIKey        string
	SpotifyClientID     string
	SpotifyClientSecret string
}

// for reference global wide
func GetConfig() Config {
	return Config{
		OpenAIAPIKey:        viper.GetString("OpenAIAPIKey"),
		SpotifyClientID:     viper.GetString("SpotifyClientID"),
		SpotifyClientSecret: viper.GetString("SpotifyClientSecret"),
	}
}

var explainFormat string = `
This program requires three inputs. then, they will be saved to "~/.local/sptfy/config.json"

1. %s: This key is used for generating playlists using GPT.
    - You can get your own key at the OpenAI API Keys page: https://platform.openai.com/account/api-keys
2. %s and %s: These are used for getting authentication with Spotify. 
    - You can create your own by creating an app at Spotify for Developers:
      https://developer.spotify.com/documentation/web-api/tutorials/getting-started#create-an-app
`

func printRequirement() {
	explain := fmt.Sprintf(explainFormat,
		color.Yellow("OPENAI_API_KEY"),
		color.Yellow("SPOTIFY_CLIENT_ID"),
		color.Yellow("SPOTIFY_CLIENT_SECRET"),
	)
	fmt.Println(explain)
}

// getConfigFromUser gets interactive input from user
func getConfigFromUser() Config {
	// prepare readline for password
	var config Config
	passwordRL, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		EnableMask:      true,
		MaskRune:        '*',
	})
	if err != nil {
		log.Fatal(err)
	}
	defer passwordRL.Close()

	// prepare prompt
	prompts := []struct {
		prompt string
		target *string
	}{
		{"Please enter OPENAI_API_KEY (input will be hidden): ", &config.OpenAIAPIKey},
		{"Please enter SPOTIFY_CLIENT_ID (input will be hidden): ", &config.SpotifyClientID},
		{"Please enter SPOTIFY_CLIENT_SECRET (input will be hidden): ", &config.SpotifyClientSecret},
	}

	// get informations
	for _, p := range prompts {
		for {
			b, err := passwordRL.ReadPassword(p.prompt)
			if err != nil {
				fmt.Println("-- tryp again. fail to read input.")
				continue
			}
			input := string(b)
			input = strings.TrimSpace(input)
			fmt.Println() // new line for aesthetics

			if input == "q" || input == "exit" {
				fmt.Println(color.Green("good bye"))
				os.Exit(0)
			}
			if err := validateInput(input); err != nil {
				fmt.Printf("-- try again. invalid input: %s\n\n", err)
				continue // invalid input, so continue the loop to re-ask
			}
			*p.target = input
			break // valid input, so break the loop and move to the next prompt
		}
	}

	return config
}

func validateInput(input string) error {
	if input == "" {
		return errors.New("input should not be empty")
	}
	return nil
}

func InitConfig() {
	// meta info
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := path.Join(homeDir, ".local", "sptfy")
	configFile := path.Join(configPath, "config.json")

	// if not exist create dir
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// if not exist set and create json
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		printRequirement()
		cfg := getConfigFromUser()
		viper.Set("OpenAIAPIKey", cfg.OpenAIAPIKey)
		viper.Set("SpotifyClientID", cfg.SpotifyClientID)
		viper.Set("SpotifyClientSecret", cfg.SpotifyClientSecret)

		if err := viper.WriteConfigAs(configFile); err != nil {
			log.Fatal(err)
		}
	}

	// now we have config.json. read config
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
