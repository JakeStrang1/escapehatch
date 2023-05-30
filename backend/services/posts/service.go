package posts

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/pages"
)

func GetPage(filter GetManyFilter, results *[]Post) (*pages.PageResult, error) {
	err := filter.Validate()
	if err != nil {
		return nil, err
	}

	page := 1
	if filter.Page != nil {
		page = *filter.Page
	}

	pageSize := defaultPerPage
	if filter.PerPage != nil {
		pageSize = *filter.PerPage
	}

	total, err := db.GetPage(filter, &Post{}, results, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &pages.PageResult{
		Page:       page,
		PerPage:    pageSize,
		TotalPages: total,
	}, nil
}

func GetByID(id string, result *Post) error {
	result.ID = id
	return db.GetByID(result)
}

func Create(userID string, result *Post) error {
	result.OwnerID = userID
	return db.Create(result)
}

func Update(id string, update PostUpdate, result *Post) error {
	result.ID = id
	err := db.GetByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(update)
	if err != nil {
		return err
	}

	return db.Update(result)
}

func Delete(id string) error {
	result := Post{}
	result.ID = id
	return db.DeleteByID(&result)
}

func GetComments(postID string, results *[]Comment) error {
	post := Post{}
	post.ID = postID
	err := db.GetByID(&post)
	if err != nil {
		return err
	}

	*results = post.Comments
	return nil
}

func GetCommentByID(postID string, commentID string, result *Comment) error {
	post := Post{}
	post.ID = postID
	err := db.GetByID(&post)
	if err != nil {
		return err
	}

	result.ID = commentID
	return post.GetCommentByID(result)
}

func CreateComment(userID string, postID string, result *Comment) error {
	post := Post{}
	post.ID = postID
	err := db.GetByID(&post)
	if err != nil {
		return err
	}

	result.OwnerID = userID
	result.PostID = postID
	err = post.AddComment(result)
	if err != nil {
		return err
	}

	return db.Update(&post)
}

func UpdateComment(postID string, commentID string, update CommentUpdate, result *Comment) error {
	post := Post{}
	post.ID = postID
	err := db.GetByID(&post)
	if err != nil {
		return err
	}

	result.ID = commentID
	err = post.GetCommentByID(result)
	if err != nil {
		return err
	}

	err = post.ApplyCommentUpdate(commentID, update, result)
	if err != nil {
		return err
	}

	return db.Update(&post)
}

func DeleteComment(postID string, commentID string) error {
	post := Post{}
	post.ID = postID
	err := db.GetByID(&post)
	if err != nil {
		return err
	}

	result := Comment{}
	result.ID = commentID
	err = post.DeleteComment(&result)
	if err != nil {
		return err
	}

	return db.Update(&post)
}
