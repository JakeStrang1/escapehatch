package posts

import (
	"fmt"

	"github.com/JakeStrang1/saas-template/db"
	"github.com/JakeStrang1/saas-template/internal/errors"
)

const defaultPerPage = 25
const maxPerPage = 250

type PostUpdate struct {
	Body *string `db:"body"`
}

type Post struct {
	db.DefaultModel `db:",inline"`
	OwnerID         string    `db:"owner_id"`
	Body            string    `db:"body"`
	Comments        []Comment `db:"comments"`
}

// BelongsTo returns true if the model belongs to the given owner
func (p *Post) BelongsTo(ownerID string) bool {
	return ownerID == p.OwnerID
}

func (p *Post) ApplyUpdate(update PostUpdate) error {
	if update.Body != nil {
		p.Body = *update.Body
	}
	return nil
}

func (p *Post) GetCommentByID(result *Comment) error {
	for _, comment := range p.Comments {
		if comment.ID == result.ID {
			*result = comment
			return nil
		}
	}

	return errors.NewNotFound()
}

func (p *Post) AddComment(result *Comment) error {
	err := result.Creating() // Call hook manually
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}

	p.Comments = append(p.Comments, *result)
	return nil
}

func (p *Post) ApplyCommentUpdate(commentID string, update CommentUpdate, result *Comment) error {
	result.ID = commentID
	err := p.GetCommentByID(result)
	if err != nil {
		return err
	}

	err = result.ApplyUpdate(update)
	if err != nil {
		return err
	}

	for i := range p.Comments {
		if p.Comments[i].ID == result.ID {
			p.Comments[i] = *result
		}
	}

	return nil
}

func (p *Post) DeleteComment(result *Comment) error {
	for i := range p.Comments {
		if p.Comments[i].ID == result.ID {
			p.Comments = append(p.Comments[:i], p.Comments[i+1:]...) // Remove element i and shift subsequent elements forward
			return nil
		}
	}
	return errors.NewNotFound()
}

type GetManyFilter struct {
	OwnerID *string `db:"owner_id,omitempty"`
	Page    *int    `db:"-"`
	PerPage *int    `db:"-"`
}

func (g *GetManyFilter) Validate() error {
	if g.Page != nil {
		if *g.Page < 1 {
			return &errors.Error{Code: errors.Invalid, Message: "page must be 1 or greater"}
		}
	}

	if g.PerPage != nil {
		if *g.PerPage < 1 || *g.PerPage > maxPerPage {
			return &errors.Error{Code: errors.Invalid, Message: fmt.Sprintf("page size must be between 1 and %d, received '%d'", maxPerPage, *g.PerPage)}
		}
	}

	return nil
}
