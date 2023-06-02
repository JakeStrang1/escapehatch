package items

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/integrations/storage"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/samber/lo"
)

/***********************************************
 * Item Services
 ***********************************************/

func Add(userID string, id string) (ItemContainer, error) {
	container, err := GetByID(id)
	if err != nil {
		return nil, err
	}

	shelfItem := users.ShelfItem{
		ItemID:      id,
		Description: container.GetItem().Description,
		ImageURL:    container.GetItem().ImageURL,
	}

	switch container.(type) {
	case *Book:
		err = users.AddBook(userID, shelfItem, &users.User{})
	case *Movie:
		err = users.AddMovie(userID, shelfItem, &users.User{})
	case *TVSeries:
		err = users.AddTVSeries(userID, shelfItem, &users.User{})
	default:
		panic("unknown item type")
	}
	if err != nil {
		return nil, err
	}

	// Get fresh item
	container, err = GetByID(id)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func GetByID(id string) (ItemContainer, error) {
	// Book
	book := newBook(id)
	if ok, err := IsBookID(&book); err != nil {
		return nil, err
	} else if ok {
		return &book, nil
	}

	// Movie
	movie := newMovie(id)
	if ok, err := IsMovieID(&movie); err != nil {
		return nil, err
	} else if ok {
		return &movie, nil
	}

	// TV Series
	tvSeries := newTVSeries(id)
	if ok, err := IsTVSeriesID(&tvSeries); err != nil {
		return nil, err
	} else if ok {
		return &tvSeries, nil
	}

	return nil, &errors.Error{Code: errors.NotFound} // ID not found under any known media type
}

func hydrateItem(item *Item) error {
	filter := users.Filter{
		ItemID: &item.ID,
	}
	count, err := users.GetCount(filter)
	if err != nil {
		return err
	}
	item.UserCount = count
	return nil
}

/***********************************************
* Book Services
***********************************************/

func CreateBook(userID string, result *Book) error {
	result.MediaType = MediaTypeBook
	result.CreatedBy = userID
	err := result.ValidateBookOnCreate()
	if err != nil {
		return err
	}

	var newImageURL string
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, storage.Options{Public: lo.ToPtr(true)})
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	err = db.Create(result)
	if err != nil {
		return err
	}

	return hydrateBook(result)
}

func IsBookID(book *Book) (bool, error) {
	err := GetBookByID(book)
	if errors.Code(err) == errors.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetBookByID(book *Book) error {
	err := db.GetByID(book)
	if err != nil {
		return err
	}
	return hydrateBook(book)
}

func hydrateBook(book *Book) error {
	book.SetDescription()
	return hydrateItem(&book.Item)
}

/***********************************************
* Movie Services
***********************************************/

func CreateMovie(userID string, result *Movie) error {
	result.MediaType = MediaTypeMovie
	result.CreatedBy = userID
	err := result.ValidateMovieOnCreate()
	if err != nil {
		return err
	}

	var newImageURL string
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, storage.Options{Public: lo.ToPtr(true)})
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	err = db.Create(result)
	if err != nil {
		return err
	}

	return hydrateMovie(result)
}

func IsMovieID(movie *Movie) (bool, error) {
	err := GetMovieByID(movie)
	if errors.Code(err) == errors.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetMovieByID(movie *Movie) error {
	err := db.GetByID(movie)
	if err != nil {
		return err
	}
	return hydrateMovie(movie)
}

func hydrateMovie(movie *Movie) error {
	movie.SetDescription()
	return hydrateItem(&movie.Item)
}

/***********************************************
* TV Series Services
***********************************************/

func CreateTVSeries(userID string, result *TVSeries) error {
	result.MediaType = MediaTypeTVSeries
	result.CreatedBy = userID
	err := result.ValidateTVSeriesOnCreate()
	if err != nil {
		return err
	}

	var newImageURL string
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, storage.Options{Public: lo.ToPtr(true)})
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	err = db.Create(result)
	if err != nil {
		return err
	}

	return hydrateTVSeries(result)
}

func IsTVSeriesID(tvSeries *TVSeries) (bool, error) {
	err := GetTVSeriesByID(tvSeries)
	if errors.Code(err) == errors.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTVSeriesByID(tvSeries *TVSeries) error {
	err := db.GetByID(tvSeries)
	if err != nil {
		return err
	}
	return hydrateTVSeries(tvSeries)
}

func hydrateTVSeries(tvSeries *TVSeries) error {
	tvSeries.SetDescription()
	return hydrateItem(&tvSeries.Item)
}
