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

// for reference global wide
func GetConfig() Config {
	return Config{
		OpenAIAPIKey:        viper.GetString("OpenAIAPIKey"),
		SpotifyClientID:     viper.GetString("SpotifyClientID"),
		SpotifyClientSecret: viper.GetString("SpotifyClientSecret"),
	}
}

// getConfigFromUser gets interactive input from user
func getConfigFromUser() Config {
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

	fmt.Println("You have to provide three config values:")
	fmt.Println("--")
	openAIAPIKey, err := passwordRL.ReadPassword("Please enter OPENAI_API_KEY: ")
	if err != nil {
		log.Fatal(err)
	}
	spotifyClientID, err := passwordRL.ReadPassword("Please enter SPOTIFY_CLIENT_ID: ")
	if err != nil {
		log.Fatal(err)
	}
	spotifyClientSecret, err := passwordRL.ReadPassword("Please enter SPOTIFY_CLIENT_SECRET: ")
	if err != nil {
		log.Fatal(err)
	}
	return Config{
		OpenAIAPIKey:        string(openAIAPIKey),
		SpotifyClientID:     string(spotifyClientID),
		SpotifyClientSecret: string(spotifyClientSecret),
	}
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
