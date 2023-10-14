package data

import (
	"time"

	"github.com/cwyang/letsgo/furthur/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // always omit
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // force string type
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

// never use space in json struct specifier
// json:",omitempty" --> omitempty without chaing key name

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 20, "title", "must not be more than 20 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greather than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive number")

	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
