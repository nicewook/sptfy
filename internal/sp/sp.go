package sp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nicewook/sptfy/internal/color"
	"github.com/nicewook/sptfy/internal/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Song struct {
	Song   string `json:"song"`
	Artist string `json:"artist"`
}

type Playlist struct {
	Playlist []Song `json:"playlist"`
}

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:9999/callback"

var (
	auth     *spotifyauth.Authenticator
	spClient *spotify.Client
	ch       = make(chan *spotify.Client)
	state    = "abc123"
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	log.Println("token:", tok)
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func AddPlaylistToSpotify(funcName string, pl Playlist) (added bool) {

	// TODO: singleton
	if auth == nil {
		log.Println("new auth for Spotify")
		auth = spotifyauth.New(
			spotifyauth.WithClientID(config.GetConfig().SpotifyClientID),
			spotifyauth.WithClientSecret(config.GetConfig().SpotifyClientSecret),
			spotifyauth.WithRedirectURL(redirectURI),
			spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPublic),
		)
	}

	// if I have access token and not expired, no need it
	if spClient != nil {
		token, err := spClient.Token()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("token: %s, token.Expiry: %T, %+v\n", token.AccessToken, token.Expiry, token.Expiry)
		log.Println("current time:", time.Now())

	} else {
		client, ok := viper.Get("spClient").(spotify.Client)
		if ok {
			log.Println("reuse client after restart!")
			spClient = &client
		}
		url := auth.AuthURL(state)

		fmt.Println(color.Hyperlink(url, color.Yellow("Click to login to Spotify")))

		// wait for auth to complete
		spClient = <-ch
		viper.Set("spClient", *spClient)
	}

	// use the client to make calls that require authorization
	user, err := spClient.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("You are logged in as:", user.ID)

	var (
		ids    []spotify.ID
		tracks []spotify.FullTrack
		ctx    = context.Background()
	)
	for _, song := range pl.Playlist {
		// https://developer.spotify.com/documentation/web-api/reference/#/operations/search
		advancdQuery := fmt.Sprintf("artist:%s track:%s", song.Song, song.Artist)
		basicQuery := fmt.Sprintf("%s %s", song.Song, song.Artist)

		for _, query := range []string{advancdQuery, basicQuery} {
			log.Println("search query:", query)
			// search for albums with the name Sempiternal
			results, err := spClient.Search(ctx, query, spotify.SearchTypeTrack)
			if err != nil {
				log.Fatal(err)
			}
			if len(results.Tracks.Tracks) == 0 {
				log.Printf("no track found for query: %s", query)
				continue
			}
			tr := results.Tracks.Tracks[0]
			if tr.Popularity < 20 {
				log.Printf("found for query: %s. but not much papular. %s, %s, popularity %d", query, tr.Name, tr.Artists, tr.Popularity, tr.PreviewURL)
				continue
			}

			log.Printf("fount %s, id: %s\n", tr.Name, tr.ID)
			ids = append(ids, tr.ID)
			tracks = append(tracks, tr)
		}
	}

	displayPlaylist(tracks)

	// create playlist

	// TODO: GPT generate playlist name and description, temperature=1
	createdPlaylist, err := spClient.CreatePlaylistForUser(
		ctx,
		user.ID,
		"my playlist by GPT",
		"generated by GPT",
		true,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	snapshotID, err := spClient.AddTracksToPlaylist(ctx, createdPlaylist.ID, ids...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("snapshotID:", snapshotID)
	return true
}

func displayPlaylist(tracks []spotify.FullTrack) {
	var trackTable [][]string
	for i, t := range tracks {
		trackTable = append(trackTable, []string{fmt.Sprintf("%02d", i+1), t.Name, t.Artists[0].Name, t.TimeDuration().String(), fmt.Sprintf("%02d", t.Popularity), t.PreviewURL})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"No", "Title", "Artist", "Duration", "Popularity", "Preview URL"})

	for _, v := range trackTable {
		table.Append(v)
	}
	table.Render()
}

func init() {
	// first start an HTTP server for OAuth
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("port is listening."))
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":9999", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
