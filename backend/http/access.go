package http

import (
	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/gin-gonic/gin"
)

type EntityType string

const (
	EntityTypeAnonymous EntityType = ""           // Zero value, a non-authenticated user
	EntityTypeUser      EntityType = "user"       // An authenticated user
	EntityTypeSuperUser EntityType = "super_user" // A user with full system access (e.g. me)
	EntityTypeSystem    EntityType = "system"     // The system itself (e.g. an automated background worker)
)

type Entity struct {
	Type EntityType
	ID   string
}

func NewAnonymousEntity() Entity {
	return Entity{
		Type: EntityTypeAnonymous,
	}
}

func NewUserEntity(userID string) Entity {
	return Entity{
		Type: EntityTypeUser,
		ID:   userID,
	}
}

func AccessPolicyAuthenticatedUsersOnly(c *gin.Context) {
	entity := c.MustGet(CtxKeyEntity).(Entity)
	if entity.Type == EntityTypeAnonymous {
		Error(c, errors.New(errors.Unauthenticated, "please sign in"))
	}
	c.Next()
}

func AccessPolicyUsersCannotOverrideSelf(c *gin.Context) {
	entity := c.MustGet(CtxKeyEntity).(Entity)
	if entity.Type == EntityTypeUser && c.GetString(CtxKeyUserID) != entity.ID {
		Error(c, errors.New(errors.Forbidden, "you cannot act on behalf of another user"))
		return
	}
	c.Next()
}

func AccessPolicyUsersCannotOverrideID(c *gin.Context) {
	entity := c.MustGet(CtxKeyEntity).(Entity)
	if entity.Type == EntityTypeUser && c.Param("id") != entity.ID && c.Param("id") != "me" {
		Error(c, errors.New(errors.Forbidden, "you cannot request this user"))
		return
	}
	c.Next()
}
