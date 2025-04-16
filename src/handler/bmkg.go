package handler

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"net/http"
)

// BMKGHandler handles BMKG related operations
type BMKGHandler struct {
	// Dependencies can be added here
}

// NewBMKGHandler creates a new instance of BMKGHandler
func NewBMKGHandler() *BMKGHandler {
	return &BMKGHandler{}
}

// AddBMKGHandler registers BMKG route handlers to the router
func (h *BMKGHandler) AddBMKGHandler(router *router.Router[*core.RequestEvent]) {
	group := router.Group("/router")
	group.GET("/hello", h.helloworld)
}

// helloworld is a simple test endpoint
func (h *BMKGHandler) helloworld(e *core.RequestEvent) error {
	return e.String(http.StatusOK, "Hello World")
}
