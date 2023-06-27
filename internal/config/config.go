package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/chzyer/readline"
	"github.com/spf13/viper"
)

type Config struct {
	OpenAIAPIKey        string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func getConfigFromUser() Config {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	fmt.Println("You have to provide three config values:")
	fmt.Println("--")
	openAIAPIKey, err := rl.ReadPassword("Please enter OPENAI_API_KEY: ")
	if err != nil {
		log.Fatal(err)
	}
	spotifyClientID, err := rl.ReadPassword("Please enter SPOTIFY_CLIENT_ID: ")
	if err != nil {
		log.Fatal(err)
	}
	spotifyClientSecret, err := rl.ReadPassword("Please enter SPOTIFY_CLIENT_SECRET: ")
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		OpenAIAPIKey:        string(openAIAPIKey),
		SpotifyClientID:     string(spotifyClientID),
		SpotifyClientSecret: string(spotifyClientSecret),
	}
}

func LoadConfig() (Config, error) {
	// meta info
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	configPath := path.Join(homeDir, ".local", "sptfy")
	configFile := path.Join(configPath, "config.json")

	// if not exist create dir
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return Config{}, err
		}
	}

	// if not exist set and create json
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := getConfigFromUser()
		viper.Set("OpenAIAPIKey", config.OpenAIAPIKey)
		viper.Set("SpotifyClientID", config.SpotifyClientID)
		viper.Set("SpotifyClientSecret", config.SpotifyClientSecret)

		if err := viper.WriteConfigAs(configFile); err != nil {
			return Config{}, err
		}
	}

	// now we have config.json. read config
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		// reset config
		return Config{}, err
	}

	config := Config{
		OpenAIAPIKey:        viper.GetString("OpenAIAPIKey"),
		SpotifyClientID:     viper.GetString("SpotifyClientID"),
		SpotifyClientSecret: viper.GetString("SpotifyClientSecret"),
	}

	return config, nil
}
