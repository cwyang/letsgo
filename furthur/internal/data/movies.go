package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // always omit
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty,string"` // force string type
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

// never use space in json struct specifier
// json:",omitempty" --> omitempty without chaing key name
