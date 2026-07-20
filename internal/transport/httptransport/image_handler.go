package httptransport

import (
	"context"
	"errors"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/image"
)

type ImageHandler struct {
	imageService *image.Service
}

func NewImageHandler(imageService *image.Service) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// Get list of user images
// (GET /api/v1/images)
func (h *ImageHandler) GetImages(
	ctx context.Context,
	r httpserver.GetImagesRequestObject,
) (httpserver.GetImagesResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.GetImages401JSONResponse{
			UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	query := &model.ImageQuery{
		Filter:     model.ImageFilter{UserID: model.IDFilter{Eq: new(user.ID)}},
		Pagination: model.Pagination{Limit: model.DefaultPaginationLimit},
	}

	if r.Params.Cursor != nil {
		cursor, err := model.DecodeCursor(*r.Params.Cursor)
		if err != nil {
			return httpserver.GetImages400JSONResponse{
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

	list, err := h.imageService.List(ctx, query)
	if err != nil {
		return httpserver.GetImages500JSONResponse{
			InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	images := make([]httpserver.Image, 0, len(list.Items))

	for _, item := range list.Items {
		images = append(images, httpserver.Image{
			ID:        item.ID,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return httpserver.GetImages200JSONResponse{
		Items: images,
		Pagination: httpserver.Pagination{
			NextCursor: list.NextCursor,
			HasNext:    list.HasNext,
		},
	}, nil
}

// Upload image file
// (POST /api/v1/images)
func (h *ImageHandler) UploadImage(
	ctx context.Context,
	r httpserver.UploadImageRequestObject,
) (httpserver.UploadImageResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.UploadImage401JSONResponse{
			UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	// todo: добавить проверку максимального размера файла (files[0].Size)

	form, err := r.Body.ReadForm(32 << 20)
	if err != nil {
		return httpserver.UploadImage400JSONResponse{
			BadRequestJSONResponse: httpserver.BadRequestJSONResponse{
				Message: "invalid multipart form",
			},
		}, nil
	}
	defer form.RemoveAll()

	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return httpserver.UploadImage400JSONResponse{
			BadRequestJSONResponse: httpserver.BadRequestJSONResponse{Message: "file is required"},
		}, nil
	}

	image, err := h.imageService.Create(ctx, user.ID, files[0])
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnsupportedImageFormat):
			return httpserver.UploadImage400JSONResponse{
				BadRequestJSONResponse: httpserver.BadRequestJSONResponse{Message: err.Error()},
			}, nil
		case errors.Is(err, model.ErrImageAlreadyExists):
			return httpserver.UploadImage400JSONResponse{
				BadRequestJSONResponse: httpserver.BadRequestJSONResponse{Message: err.Error()},
			}, nil
		default:
			return httpserver.UploadImage500JSONResponse{
				InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
					Message: "internal server error",
				},
			}, nil
		}
	}

	return httpserver.UploadImage201JSONResponse{
		ID:        image.ID,
		Name:      image.Name,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}, nil
}

// Delete
// (DELETE /api/v1/images/{name})
func (h *ImageHandler) DeleteImageByName(
	ctx context.Context,
	r httpserver.DeleteImageByNameRequestObject,
) (httpserver.DeleteImageByNameResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.DeleteImageByName401JSONResponse{
			UnauthorizedJSONResponse: httpserver.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	if err := h.imageService.Delete(ctx, user.ID, r.Name); err != nil {
		return httpserver.DeleteImageByName500JSONResponse{
			InternalErrorJSONResponse: httpserver.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return httpserver.DeleteImageByName204Response{}, nil
}
