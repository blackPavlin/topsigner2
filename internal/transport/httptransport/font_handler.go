package httptransport

import (
	"context"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
)

type FontHandler struct {
	fontService *font.Service
}

func NewFontHandler(fontService *font.Service) *FontHandler {
	return &FontHandler{
		fontService: fontService,
	}
}

// Get list of fonts
// (GET /api/v1/fonts)
func (h *FontHandler) GetFonts(
	ctx context.Context,
	r httpserver.GetFontsRequestObject,
) (httpserver.GetFontsResponseObject, error) {
	if _, ok := auth.GetUserFromContext(ctx); !ok {
		return httpserver.GetFonts401JSONResponse{
			UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	query := &model.FontQuery{
		Pagination: model.Pagination{Limit: model.DefaultPaginationLimit},
	}

	if r.Params.Cursor != nil {
		cursor, err := model.DecodeCursor(*r.Params.Cursor)
		if err != nil {
			return httpserver.GetFonts400JSONResponse{
				BadRequestJSONResponse: httpserver.BadRequestJSONResponse{Message: err.Error()},
			}, nil
		}

		query.Pagination.Cursor = cursor
	}

	if r.Params.Limit != nil {
		query.Pagination.Limit = *r.Params.Limit

		if query.Pagination.Limit <= 0 || query.Pagination.Limit > model.DefaultPaginationLimit {
			query.Pagination.Limit = model.DefaultPaginationLimit
		}
	}

	list, err := h.fontService.List(ctx, query)
	if err != nil {
		return httpserver.GetFonts500JSONResponse{
			InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	fonts := make([]httpserver.Font, 0, len(list.Items))

	for _, item := range list.Items {
		fonts = append(fonts, httpserver.Font{
			ID:        item.ID,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return httpserver.GetFonts200JSONResponse{
		Items: fonts,
		Pagination: httpserver.Pagination{
			NextCursor: list.NextCursor,
			HasNext:    list.HasNext,
		},
	}, nil
}
