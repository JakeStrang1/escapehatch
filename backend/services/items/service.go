package items

import (
	"fmt"

	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/integrations/storage"
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/internal/pages"
	"github.com/JakeStrang1/escapehatch/services/users"
	"github.com/kamva/mgm/v3"
	"github.com/samber/lo"
)

const defaultPerPage = 25
const maxPerPage = 250

/***********************************************
 * Item Services
 ***********************************************/

func GetPage(filter Filter) ([]any, *pages.PageResult, error) {
	err := filter.Validate()
	if err != nil {
		return nil, nil, err
	}

	page := 1
	if filter.Page != nil {
		page = *filter.Page
	}

	pageSize := defaultPerPage
	if filter.PerPage != nil {
		pageSize = *filter.PerPage
	}

	opts := db.GetManyOptions{Sort: [][]any{{"user_count", -1}}} // Sort by most popular
	itemStats := []ItemStat{}
	total, err := db.GetPage(filter, &ItemStat{}, &itemStats, page, pageSize, opts)
	if err != nil {
		return nil, nil, err
	}

	results := []any{}
	for _, itemStat := range itemStats {
		switch itemStat.MediaType {
		case MediaTypeBook:
			book := newBook(itemStat.ItemID)
			err = GetBookByID(&book)
			if err != nil {
				return nil, nil, err
			}
			results = append(results, &book)
		case MediaTypeMovie:
			movie := newMovie(itemStat.ItemID)
			err = GetMovieByID(&movie)
			if err != nil {
				return nil, nil, err
			}
			results = append(results, &movie)
		case MediaTypeTVSeries:
			tvSeries := newTVSeries(itemStat.ItemID)
			err = GetTVSeriesByID(&tvSeries)
			if err != nil {
				return nil, nil, err
			}
			results = append(results, &tvSeries)
		default:
			panic("unknown type")
		}
	}

	return results, &pages.PageResult{
		Page:       page,
		PerPage:    pageSize,
		TotalPages: total,
	}, nil
}

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

	// Save stat
	fmt.Println(id)
	err = db.IncrementOne(db.M{"item_id": id}, "user_count", &ItemStat{})
	if err != nil {
		return nil, err
	}

	return container, nil
}

func Remove(userID string, id string) (ItemContainer, error) {
	err := users.RemoveItemFromAllShelves(userID, id, &users.User{})
	if err != nil {
		return nil, err
	}

	// Get item
	container, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	return container, nil
}

type DeleteParams struct {
	UserID string
	Reason string
	ItemID string
}

func (d *DeleteParams) Validate() error {
	if d.Reason == "" {
		return &errors.Error{Code: errors.Invalid, Message: "reason must not be blank"}
	}
	return nil
}

func Delete(params DeleteParams) error {
	// Validate
	err := params.Validate()
	if err != nil {
		return err
	}

	// Get item
	container, err := GetByID(params.ItemID)
	if err != nil {
		return err
	}

	filter := users.Filter{
		ItemID: &params.ItemID,
	}
	count, err := users.GetCount(filter)
	if err != nil {
		return err
	}

	// Backup in deleted_items
	_, err = mgm.CollectionByName("deleted_items").InsertOne(mgm.Ctx(), container)
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}

	// Delete
	model := container.(mgm.Model)
	err = db.DeleteByID(model)
	if err != nil {
		return err
	}

	// Track Changes
	Track(container.GetItem().ID).Deleted(params.Reason, count).By(params.UserID).Save()

	// Save stat
	db.DeleteOne(db.M{"item_id": params.ItemID}, &ItemStat{})

	// Remove from shelves
	err = users.RemoveItemFromAllUsers(params.ItemID)
	if err != nil {
		return err
	}
	return nil
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

func SaveImage(result *Item) error {
	var newImageURL string
	var err error
	options := storage.Options{
		Public:         lo.ToPtr(true),
		ImageCompress:  lo.ToPtr(true),
		ImageMaxWidth:  lo.ToPtr(600),
		ImageMaxHeight: lo.ToPtr(600),
		ImageMaxKB:     lo.ToPtr(100),
	}
	if len(result.ImageFileBody) != 0 {
		newImageURL, err = storage.Create(result.ImageFileName, result.ImageFileBody, options)
	} else {
		newImageURL, err = storage.CreateFromURL(result.ImageURL, options)
	}
	if err != nil {
		return err
	}
	result.ImageURL = newImageURL
	return nil
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
	err := result.ValidateOnCreate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.EnsureTextIndex(result, BookSearchPaths)
	if err != nil {
		return err
	}

	err = db.Create(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(result.ID).Created(result).By(userID).Save()

	// Save Stat
	db.Create(result.ToStat(MediaTypeBook))

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

func UpdateBook(userID string, id string, update BookUpdate, result *Book) error {
	result.ID = id
	err := GetBookByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(userID, update)
	if err != nil {
		return err
	}

	err = result.Validate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.Update(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(id).Updated(update, result).By(userID).Save()

	// TODO: Refresh cached image links and descriptions on user shelves (Cloud Tasks looks like good option)

	return hydrateBook(result)
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
	err := result.ValidateOnCreate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.EnsureTextIndex(result, MovieSearchPaths)
	if err != nil {
		return err
	}

	err = db.Create(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(result.ID).Created(result).By(userID).Save()

	// Save Stat
	db.Create(result.ToStat(MediaTypeMovie))

	return hydrateMovie(result)
}

func UpdateMovie(userID string, id string, update MovieUpdate, result *Movie) error {
	result.ID = id
	err := GetMovieByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(userID, update)
	if err != nil {
		return err
	}

	err = result.Validate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.Update(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(id).Updated(update, result).By(userID).Save()

	// TODO: Refresh cached image links and descriptions on user shelves (Cloud Tasks looks like good option)

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
	err := result.ValidateOnCreate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.EnsureTextIndex(result, TVSeriesSearchPaths)
	if err != nil {
		return err
	}

	err = db.Create(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(result.ID).Created(result).By(userID).Save()

	// Save Stat
	db.Create(result.ToStat(MediaTypeTVSeries))

	return hydrateTVSeries(result)
}

func UpdateTVSeries(userID string, id string, update TVSeriesUpdate, result *TVSeries) error {
	result.ID = id
	err := GetTVSeriesByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(userID, update)
	if err != nil {
		return err
	}

	err = result.Validate()
	if err != nil {
		return err
	}

	err = SaveImage(&result.Item)
	if err != nil {
		return err
	}

	err = db.Update(result)
	if err != nil {
		return err
	}

	// Track Changes
	Track(id).Updated(update, result).By(userID).Save()

	// TODO: Refresh cached image links and descriptions on user shelves (Cloud Tasks looks like good option)

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
