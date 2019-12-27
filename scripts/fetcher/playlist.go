package main

import (
	"github.com/zmb3/spotify"
	"regexp"
)

// Custom types to create an easy data structure for Hugo data
type Playlist struct {
	Name        string
	Description string
	Url         string
	Type        string
	Collection  string
	DateStart   string
	DateLast    string
	Image       []spotify.Image
}

// Extract data from a Spotifakt compatible playlist name
func ParseName(n string, pattern string, def string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(n)
	if match == nil {
		return def
	}
	return string(match[1])
}

// Create a new playlist
func NewPlaylist(p *spotify.FullPlaylist) *Playlist {
	obj := new(Playlist)
	obj.Description = p.Description
	obj.Url = p.ExternalURLs["spotify"]
	obj.Name = ParseName(p.Name, `(?U:(.+)(?:\s\())`, "")
	obj.Type = ParseName(p.Name, `(?U:(.+)(?:\s\|\s))`, "")
	obj.Collection = ParseName(p.Name, `(?U:(?:\s\|\s)([\sa-zA-Z]+):)`, "Misc")
	obj.DateStart = "2019"
	obj.DateLast = "2019"
	obj.Image = p.Images
	return obj
}
