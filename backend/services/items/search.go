package items

import (
	"sort"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal"
	"github.com/JakeStrang1/escapehatch/services/users"
)

type SearchResult interface {
	Score() float64
	Result() any
}

type UserSearch struct {
	users.User  `db:",inline"`
	SearchScore float64 `db:"search_score"`
}

func (u *UserSearch) Score() float64 {
	return u.SearchScore
}

func (u *UserSearch) Result() any {
	return u.User
}

type BookSearch struct {
	Book        `db:",inline"`
	SearchScore float64 `db:"search_score"`
}

func (u *BookSearch) Score() float64 {
	return u.SearchScore
}

func (u *BookSearch) Result() any {
	return u.Book
}

type MovieSearch struct {
	Movie       `db:",inline"`
	SearchScore float64 `db:"search_score"`
}

func (u *MovieSearch) Score() float64 {
	return u.SearchScore
}

func (u *MovieSearch) Result() any {
	return u.Movie
}

type TVSeriesSearch struct {
	TVSeries    `db:",inline"`
	SearchScore float64 `db:"search_score"`
}

func (u *TVSeriesSearch) Score() float64 {
	return u.SearchScore
}

func (u *TVSeriesSearch) Result() any {
	return u.TVSeries
}

func GetSearch(search string, results *[]any) error {
	searchResults := []SearchResult{}

	// Users
	userSearchResults := []UserSearch{}
	err := db.Search(search, users.UserSearchPaths, &users.User{}, &userSearchResults)
	if err != nil {
		return err
	}
	searchResults = append(searchResults, internal.Map(userSearchResults, func(r UserSearch) SearchResult { return &r })...)

	// Books
	bookSearchResults := []BookSearch{}
	err = db.Search(search, BookSearchPaths, &Book{}, &bookSearchResults)
	if err != nil {
		return err
	}
	searchResults = append(searchResults, internal.Map(bookSearchResults, func(r BookSearch) SearchResult { return &r })...)

	// Movies
	movieSearchResults := []MovieSearch{}
	err = db.Search(search, MovieSearchPaths, &Movie{}, &movieSearchResults)
	if err != nil {
		return err
	}
	searchResults = append(searchResults, internal.Map(movieSearchResults, func(r MovieSearch) SearchResult { return &r })...)

	// TV series
	tvSeriesSearchResults := []TVSeriesSearch{}
	err = db.Search(search, TVSeriesSearchPaths, &TVSeries{}, &tvSeriesSearchResults)
	if err != nil {
		return err
	}
	searchResults = append(searchResults, internal.Map(tvSeriesSearchResults, func(r TVSeriesSearch) SearchResult { return &r })...)

	// Sort
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Score() > searchResults[j].Score() // return higher scores first
	})

	*results = internal.Map(searchResults, func(r SearchResult) any { return r.Result() }) // Might want a hydrate here?
	return nil
}
