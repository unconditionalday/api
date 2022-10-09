package app

import "time"

type Feed struct {
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Language string    `json:"language"`
	Image    *Image    `json:"image"`
	Summary  string    `json:"summary"`
	Source   string    `json:"source"`
	Date     time.Time `json:"date"`
}

type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
