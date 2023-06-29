package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/nicewook/sptfy/internal/ai"
	"github.com/nicewook/sptfy/internal/color"
	"github.com/nicewook/sptfy/internal/sp"
)

func init() {
	runMode := os.Getenv("RUN_MODE")
	if runMode != "dev" {
		log.SetOutput(io.Discard)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("sptfy start")

	rl := getReadline()
	defer rl.Close()
	rl.CaptureExitSignal()

	for {
		fmt.Println()
		prompt, num := getPrompt(rl)
		fmt.Println()

		funcName, playlist := ai.GeneratePlaylist(prompt, num)
		fmt.Println()

		if funcName == "" || len(playlist.Playlist) == 0 {
			fmt.Println("try again. fail to genereate playlist.")
			continue
		}
		log.Println("function called:", funcName)

		b, err := json.MarshalIndent(playlist, "  ", "")
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		playlistName := ai.GeneratePlayListName(prompt, string(b))
		fmt.Println("Playlist name:", playlistName)
		fmt.Println()

		added := sp.AddPlaylistToSpotify(playlistName, playlist)
		if !added {
			fmt.Println("failed to add to Spotify.")
		}
	}
}

func getPrompt(rl *readline.Instance) (string, int) {
	fmt.Println(color.Blue("Let's make a playlist to Spotify!"))
	fmt.Println()
	fmt.Printf("Describe music playlist you want to listen(or %s):\n", color.Yellow("exit, q"))
	prompt, err := rl.Readline()
	if err == readline.ErrInterrupt {
		if len(prompt) == 0 {
			log.Fatal("interrupted")
		}
	} else if err != nil {
		log.Fatal(err)
	}
	prompt = strings.TrimSpace(prompt)
	if prompt == "exit" || prompt == "q" {
		fmt.Println(color.Green("good bye"))
		os.Exit(0)
	}
	fmt.Println("How many songs? (4 to 20, can be generated less:")

	numStr, err := rl.Readline()
	if err == readline.ErrInterrupt {
		if len(prompt) == 0 {
			log.Fatal("interrupted")
		}
	} else if err != nil {
		log.Fatal(err)
	}
	numStr = strings.TrimSpace(numStr)
	if numStr == "exit" || numStr == "q" {
		fmt.Println(color.Green("good bye"))
		os.Exit(0)
	}
	num, err := strconv.Atoi(numStr)
	if err != nil || (num < 4 && num > 20) {
		fmt.Println(color.Red("try again. it should be a number from 4 to 20"))
	}
	return prompt, num
}
