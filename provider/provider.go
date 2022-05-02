package provider

import (
	"github.com/javtube/javtube-sdk-go/model"
)

var (
	// UA is the default user agent for each provider.
	UA = colly.UserAgent(UA)()
)

type Provider interface {
	// GetMovieInfoByID gets movie's info by id.
	GetMovieInfoByID(id string) (*model.MovieInfo, error)

	// GetMovieInfoByLink gets movie's info by link.
	GetMovieInfoByLink(link string) (*model.MovieInfo, error)

	// SearchMovie searches matched movies.
	SearchMovie(keyword string) ([]*model.SearchResult, error)
}
