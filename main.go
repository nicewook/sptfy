package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

		prompt := getPrompt(rl)
		fmt.Println()

		funcName, playlist := ai.GeneratePlaylist(prompt)
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

const (
	q01 = "Q: What is the primary purpose of this playlist? (For example, workout, study, party, relaxation, driving, commute, reading.)"

	q02 = "Q: Can you name a few artists or bands that you would like to include in this playlist? We will also add musicians of similar taste."

	q03 = "Q: Are there any specific songs you would like me to add to the playlist? We will also add songs of similar taste."

	q04 = "Q: What genre or style of music do you prefer? (For example, pop, rock, country, classical, hip-hop, jazz.)"

	q05 = "Q: Do you have any preferences for a particular era or time period of music? (For example, 80s, 90s, modern.)"

	q06 = "Q: What kind of mood would you like the playlist to convey? (For example, happy, sad, energetic, chill, romantic.)"

	q07 = "Q: Would you like the playlist to include songs with lyrics, instrumental tracks, or a mix of both?"

	q08 = "Q: Do you have any preference for the language or contury of the songs?"

	q09 = "Q: How long do you want the playlist to be? (For example, a certain number of songs, or a certain duration in hours or minutes)"

	q10 = "Q: Are you open to discovering new artists or songs that are similar to your preferences?"

	q11 = "Q: Would you like a variety of tempos in the playlist or would you prefer a consistent tempo?"
)

// exmample
// Q: What is the primary purpose of this playlist? (For example, workout, study, party, relaxation, driving, commute.)A: commte
// Q: Can you name a few artists or bands that you would like to include in this playlist?A: fabrizio paterlini
// Q: Are there any specific songs you would like me to add to the playlist?A:
// Q: What genre or style of music do you prefer? (For example, pop, rock, country, classical, hip-hop, jazz.)A: jazz
// Q: Do you have any preferences for a particular era or time period of music? (For example, 80s, 90s, modern.)A:
// Q: What kind of mood would you like the playlist to convey? (For example, happy, sad, energetic, chill, romantic.)A: chill
// Q: Would you like the playlist to include songs with lyrics, instrumental tracks, or a mix of both?A: instrumental
// Q: Do you have any preference for the language or contury of the songs?A: nope
// Q: How long do you want the playlist to be? (For example, a certain number of songs, or a certain duration in hours or minutes)A: 20 min
// Q: Are you open to discovering new artists or songs that are similar to your preferences?A: yes
// Q: Would you like a variety of tempos in the playlist or would you prefer a consistent tempo?A: consistent

func getPrompt(rl *readline.Instance) string {
	var result string
	fmt.Println(color.Blue("Let's make a playlist to Spotify!"))
	fmt.Printf("Answer following questions to generate a better playlist of your taste(or %s to quit):\n", color.Yellow("exit, q"))
	fmt.Printf("Also, you can skip the question by %s\n", "pressing enter or type skip")

	getAnswer(rl, q01, &result)
	getAnswer(rl, q02, &result)
	getAnswer(rl, q03, &result)
	getAnswer(rl, q04, &result)
	getAnswer(rl, q05, &result)
	getAnswer(rl, q06, &result)
	getAnswer(rl, q07, &result)
	getAnswer(rl, q08, &result)
	getAnswer(rl, q09, &result)
	getAnswer(rl, q10, &result)
	getAnswer(rl, q11, &result)

	log.Println("final prompt:", result)

	return result
}

func getAnswer(rl *readline.Instance, question string, result *string) {
	fmt.Println(question)
	input := getInput(rl)
	log.Println("input:", input)
	if input != "" && input != "skip" {
		// *result += fmt.Sprintln(question)
		*result += question
		*result += fmt.Sprintln("A:", input)
	}
}

func getInput(rl *readline.Instance) string {
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
	return prompt
}
