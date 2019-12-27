// Fetch public playlists with metadata on a provided spotify account.
//
// Spotify API credentials should be exported in environment or .env file
// See .env.dist for specifications
//
// Usage:
//   go run main.go --user=<your spotify user id>
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
)

var userID = flag.String("user", "_u.zer_", "the Spotify user ID to look up")

func main() {
	flag.Parse()
	if *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: missing user ID\n")
		flag.Usage()
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("Couldn't get token: %v", err)
	}
	client := spotify.Authenticator{}.NewClient(token)
	user, err := client.GetUsersPublicProfile(spotify.ID(*userID))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	// Get playlists
	pages, err := client.GetPlaylistsForUser(user.ID)
	if err != nil {
		log.Fatalf("Couldn't get features playlists: %v", err)
	}

	// Print results
	for page := 1; ; page++ {
		for _, pl := range pages.Playlists {
			fields := "name,description,followers,images,external_urls,tracks.total"
			playlist, err := client.GetPlaylistOpt(pl.ID, fields)
			if err != nil {
				log.Fatalf("Couldn't get full playlist metadata: %v", err)
			}

			// Output a filtered JSON
			b, err := json.Marshal(playlist)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}
		err = client.NextPage(pages)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
