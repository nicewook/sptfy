package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chzyer/readline"
	"github.com/nicewook/sptfy/internal/ai"
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

	fmt.Println()
	fmt.Println("Describe music playlist you want to listen:")
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
		funcName, playlist := ai.GeneratePlaylist(prompt, 8)
		if funcName == "" || len(playlist.Playlist) == 0 {
			fmt.Println("try again. fail to genereate playlist.")
			continue
		}
		// TODO
		// 1. display as table
		// 2. ask if you want to add
		added := sp.AddPlaylistToSpotify(funcName, playlist)
		if added {
			fmt.Println("successfully added to Spotify.")
			// TODO: spotify URL link
			continue
		}
		fmt.Println("failed to add to Spotify.")
	}
}
