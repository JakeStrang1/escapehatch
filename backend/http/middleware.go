package http

import (
	"net/http"

	"github.com/JakeStrang1/escapehatch/internal/errors"
	"github.com/JakeStrang1/escapehatch/services/auth"
	"github.com/gin-gonic/gin"
)

const CtxKeyEntity = "entity"
const CtxKeyUserID = "userID"
const UserIDQuery = "userID" // An optional query param that can be provided to any request to override the user

func Authenticate(c *gin.Context) {
	sessionToken, err := c.Cookie("SID")
	if errors.Is(err, http.ErrNoCookie) {
		// Error(c, &errors.Error{Code: errors.Unauthenticated, Message: "please sign in", Err: err}) // This is the old response
		c.Set(CtxKeyEntity, NewAnonymousEntity())
		return
	}
	if err != nil {
		// I don't think this is reachable
		Error(c, &errors.Error{Code: errors.Internal, Err: err})
		return
	}

	// Authenticate user
	user, err := auth.Authenticate(sessionToken)
	if err != nil {
		Error(c, err)
		return
	}
	c.Set(CtxKeyEntity, NewUserEntity(user.ID))

	// The userID context key tells handlers who to execute the request on behalf of.
	// This is not necessarily the same as the calling entity.
	// Access control must make sure that the entity and userID combination are compatible with the request.
	if userID, ok := c.GetQuery(UserIDQuery); ok {
		c.Set(CtxKeyUserID, userID)
	} else {
		c.Set(CtxKeyUserID, user.ID)
	}

	c.Next()
}
