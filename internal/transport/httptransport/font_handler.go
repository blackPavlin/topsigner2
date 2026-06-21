package httptransport

import (
	"context"

	"github.com/bboykiv/topsigner/gen/openapi"
)

type FontHandler struct{}

func NewFontHandler() *FontHandler {
	return &FontHandler{}
}

// Get list of fonts
// (GET /api/v1/fonts)
func (h *FontHandler) GetFonts(
	ctx context.Context,
	r openapi.GetFontsRequestObject,
) (openapi.GetFontsResponseObject, error) {
	return nil, nil
}
