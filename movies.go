package main

import (
	"fmt"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/pkg/errors"
)

const (
	locale             = "en-US"
	posterThumbTmpl    = "http://image.tmdb.org/t/p/w92/%s"
	posterOriginalTmpl = "http://image.tmdb.org/t/p/original/%s"
)

type Movie struct {
	ID            string
	Title         string
	OriginalTitle string
	PosterImg     string
	PosterThumb   string
	Overview      string
	ReleaseDate   string
	VoteCount     int64
	Popularity    float32
	VoteAvg       float32
}

func searchMovies(tmdbAPI *tmdb.Client, input string, offset string) ([]Movie, string, error) {
	tmdbClient, err := tmdb.Init(conf.TmdbAPIKey)
	if err != nil {
		return nil, "", err
	}

	result, err := tmdbClient.GetSearchMovies(input, map[string]string{
		"language": locale,
		"page":     offset,
	})

	if err != nil {
		return nil, "", errors.Wrap(err, "tmdb movies search")
	}

	return mapToMovies(result), strconv.Itoa(int(result.Page + 1)), nil
}

func mapToMovies(result *tmdb.SearchMovies) []Movie {
	movies := make([]Movie, 0, len(result.Results))
	for _, r := range result.Results {
		movies = append(movies, Movie{
			ID:            strconv.Itoa(int(r.ID)),
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			Overview:      r.Overview,
			ReleaseDate:   r.ReleaseDate,
			PosterImg:     fmt.Sprintf(posterOriginalTmpl, r.PosterPath),
			PosterThumb:   fmt.Sprintf(posterThumbTmpl, r.PosterPath),
			Popularity:    r.Popularity,
			VoteAvg:       r.VoteAverage,
			VoteCount:     r.VoteCount,
		})
	}

	return movies
}
