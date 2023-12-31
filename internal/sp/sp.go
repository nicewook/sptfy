package sp

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nicewook/sptfy/internal/color"
	"github.com/nicewook/sptfy/internal/config"
	tw "github.com/olekukonko/tablewriter"
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
	state    string
)

func init() {
	state = generateState()

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

func generateState() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

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
	log.Println("token expires:", tok.Expiry.Format(time.RFC3339))
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func AddPlaylistToSpotify(playlistName string, pl Playlist) (added bool) {

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

	// need to login
	if spClient == nil {
		url := auth.AuthURL(state)
		fmt.Println(color.Blue("Create a playlist on Spotify"))
		fmt.Println(color.Hyperlink(url, color.Yellow("Click to authenticate to Spotify!")))
		spClient = <-ch // wait for auth to complete
	}

	// just debugging
	token, err := spClient.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("token: %s, token.Expiry: %T, %+v\n", token.AccessToken, token.Expiry, token.Expiry)

	// use the client to make calls that require authorization
	user, err := spClient.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("logged in user id:", user.ID)

	var (
		ids    []spotify.ID
		tracks []spotify.FullTrack
		ctx    = context.Background()
	)
	for _, song := range pl.Playlist {
		// https://developer.spotify.com/documentation/web-api/reference/#/operations/search
		basicQuery := fmt.Sprintf("%s %s", song.Song, song.Artist)

		log.Println("search query:", basicQuery)
		results, err := spClient.Search(ctx, basicQuery, spotify.SearchTypeTrack|spotify.SearchTypeArtist)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(results.Tracks.Tracks) == 0 {
			log.Println("no track found for query:", basicQuery)
			continue
		}
		tr := results.Tracks.Tracks[0]
		if tr.Popularity < 20 {
			log.Printf("found for query: %s. but not much papular. %d", basicQuery, tr.Popularity)
			continue
		}

		log.Printf("fount %s, id: %s\n", tr.Name, tr.ID)
		ids = append(ids, tr.ID)
		tracks = append(tracks, tr)
	}

	// create playlist

	// TODO: GPT generate playlist name and description, temperature=1
	plDescription := fmt.Sprintf("generated by sptfy app using GPT on %v", time.Now().Format("2006-01-02 15:04:05 MST"))
	createdPlaylist, err := spClient.CreatePlaylistForUser(
		ctx,
		user.ID,
		playlistName,
		plDescription,
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
	log.Println("snapshotID: ", snapshotID)
	if url, exist := createdPlaylist.ExternalURLs["spotify"]; exist {
		fmt.Println(color.Blue("Successfully created on Spotify!"))
		fmt.Println(color.Hyperlink(url, color.Yellow(playlistName)))
		displayPlaylist(tracks)
	}
	return true
}

func displayPlaylist(tracks []spotify.FullTrack) {
	var trackTable [][]string
	for i, t := range tracks {
		trackTable = append(trackTable, []string{fmt.Sprintf("%02d", i+1), t.Name, t.Artists[0].Name, formatPlaytime(t.TimeDuration()), fmt.Sprintf("%02d", t.Popularity), t.PreviewURL})
	}

	table := tw.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"No", "Title", "Artist", "Play time", "Popularity", "Preview URL"})
	table.SetAutoFormatHeaders(false)
	table.SetColumnAlignment([]int{tw.ALIGN_CENTER, tw.ALIGN_LEFT, tw.ALIGN_LEFT, tw.ALIGN_RIGHT, tw.ALIGN_CENTER, tw.ALIGN_LEFT})

	for _, v := range trackTable {
		table.Append(v)
	}
	table.Render()
}

func formatPlaytime(playtime time.Duration) string {
	minutes := int(playtime.Minutes())
	seconds := int(playtime.Seconds()) - (minutes * 60)

	return fmt.Sprintf("%01d:%02d", int(minutes), int(seconds))
}
