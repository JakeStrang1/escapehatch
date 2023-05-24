package posts

import (
	"github.com/JakeStrang1/escapehatch/db"
	"github.com/JakeStrang1/escapehatch/internal/errors"
)

type Comment struct {
	db.DefaultModel `db:",inline"`
	OwnerID         string `db:"owner_id"`
	PostID          string `db:"post_id"`
	Body            string `db:"body"`
	ReplyTo         string `db:"reply_to"` // A comment ID of a separate comment on this post
}

type CommentUpdate struct {
	Body *string `db:"body"`
}

func (s *Comment) ApplyUpdate(update CommentUpdate) error {
	if update.Body != nil {
		s.Body = *update.Body
	}

	err := s.Saving() // Call hook manually
	if err != nil {
		return &errors.Error{Code: errors.Internal, Err: err}
	}

	return nil
}
