package httptransport

import (
	"context"
	"errors"

	"github.com/bboykiv/topsigner/gen/openapi"
	"github.com/bboykiv/topsigner/internal/model"
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
	r openapi.GetImagesRequestObject,
) (openapi.GetImagesResponseObject, error) {
	user, ok := model.GetUserFromContext(ctx)
	if !ok {
		return openapi.GetImages401JSONResponse{
			UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	query := &model.ImageQuery{
		Filter:     model.ImageFilter{UserID: model.IDFilter{Eq: new(user.ID)}},
		Pagination: model.Pagination{Limit: model.DefaultPaginationLimit},
	}

	if r.Params.Cursor != nil {
		cursor, err := model.DecodeCursor(*r.Params.Cursor)
		if err != nil {
			return openapi.GetImages400JSONResponse{
				BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: err.Error()},
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
		return openapi.GetImages500JSONResponse{
			InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	images := make([]openapi.Image, 0, len(list.Items))

	for _, item := range list.Items {
		images = append(images, openapi.Image{
			ID:        item.ID,
			Name:      item.Name,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return openapi.GetImages200JSONResponse{
		Items: images,
		Pagination: openapi.Pagination{
			NextCursor: list.NextCursor,
			HasNext:    list.HasNext,
		},
	}, nil
}

// Upload image file
// (POST /api/v1/images)
func (h *ImageHandler) UploadImage(
	ctx context.Context,
	r openapi.UploadImageRequestObject,
) (openapi.UploadImageResponseObject, error) {
	user, ok := model.GetUserFromContext(ctx)
	if !ok {
		return openapi.UploadImage401JSONResponse{
			UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	// todo: добавить проверку максимального размера файла (files[0].Size)

	form, err := r.Body.ReadForm(32 << 20)
	if err != nil {
		return openapi.UploadImage400JSONResponse{
			BadRequestJSONResponse: openapi.BadRequestJSONResponse{
				Message: "invalid multipart form",
			},
		}, nil
	}
	defer form.RemoveAll()

	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return openapi.UploadImage400JSONResponse{
			BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: "file is required"},
		}, nil
	}

	image, err := h.imageService.Create(ctx, user.ID, files[0])
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnsupportedImageFormat):
			return openapi.UploadImage400JSONResponse{
				BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: err.Error()},
			}, nil
		case errors.Is(err, model.ErrImageAlreadyExists):
			return openapi.UploadImage400JSONResponse{
				BadRequestJSONResponse: openapi.BadRequestJSONResponse{Message: err.Error()},
			}, nil
		default:
			return openapi.UploadImage500JSONResponse{
				InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
					Message: "internal server error",
				},
			}, nil
		}
	}

	return openapi.UploadImage201JSONResponse{
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
	r openapi.DeleteImageByNameRequestObject,
) (openapi.DeleteImageByNameResponseObject, error) {
	user, ok := model.GetUserFromContext(ctx)
	if !ok {
		return openapi.DeleteImageByName401JSONResponse{
			UnauthorizedJSONResponse: openapi.UnauthorizedJSONResponse{Message: "unauthorized"},
		}, nil
	}

	if err := h.imageService.Delete(ctx, user.ID, r.Name); err != nil {
		return openapi.DeleteImageByName500JSONResponse{
			InternalErrorJSONResponse: openapi.InternalErrorJSONResponse{
				Message: "internal server error",
			},
		}, nil
	}

	return openapi.DeleteImageByName204Response{}, nil
}
