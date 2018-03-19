package routes

import (
	"github.com/go-kit/kit/log"

	"github.com/coyle/bridge/server/routes/users"
)

// Handler contains all route handlers for a service
type Handler struct {
	Logger log.Logger
	User   *users.User
}
