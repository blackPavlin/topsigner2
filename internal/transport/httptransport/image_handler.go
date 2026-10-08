package httptransport

import (
	"context"
	"errors"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/transport/httptransport/middleware"
)

type ImageHandler struct {
	imageService *image.Service
}

func NewImageHandler(imageService *image.Service) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// ListImages Get list of user images
// (GET /api/v1/images)
func (h *ImageHandler) ListImages(
	ctx context.Context,
	r httpserver.ListImagesRequestObject,
) (httpserver.ListImagesResponseObject, error) {
	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return httpserver.ListImages401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	query := &model.ImageQuery{
		Filter:     model.ImageFilter{UserID: model.IDFilter{Eq: new(user.ID)}},
		Pagination: model.Pagination{Limit: model.DefaultPaginationLimit},
	}

	if r.Params.Cursor != nil {
		cursor, err := model.DecodeCursor(*r.Params.Cursor)
		if err != nil {
			return httpserver.ListImages400JSONResponse{
				BadRequestJSONResponse: NewBadRequestError(err),
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
		return httpserver.ListImages500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
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

	return httpserver.ListImages200JSONResponse{
		Items: images,
		Pagination: httpserver.Pagination{
			NextCursor: list.NextCursor,
		},
	}, nil
}

// UploadImage Upload image file
// (POST /api/v1/images)
func (h *ImageHandler) UploadImage(
	ctx context.Context,
	r httpserver.UploadImageRequestObject,
) (httpserver.UploadImageResponseObject, error) {
	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return httpserver.UploadImage401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	// todo: добавить проверку максимального размера файла (files[0].Size)

	form, err := r.Body.ReadForm(32 << 20)
	if err != nil {
		return httpserver.UploadImage400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(ErrInvalidMultipartForm),
		}, nil
	}
	defer form.RemoveAll()

	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return httpserver.UploadImage400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(ErrMultipartFileIsRequired),
		}, nil
	}

	file, err := files[0].Open()
	if err != nil {
		return httpserver.UploadImage400JSONResponse{
			BadRequestJSONResponse: NewBadRequestError(ErrInvalidMultipartForm),
		}, nil
	}
	defer file.Close()

	image, err := h.imageService.Create(ctx, user.ID, &image.Upload{
		Filename: files[0].Filename,
		Size:     files[0].Size,
		File:     file,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnsupportedImageFormat):
			return httpserver.UploadImage400JSONResponse{
				BadRequestJSONResponse: NewBadRequestError(err),
			}, nil
		case errors.Is(err, model.ErrImageAlreadyExists):
			return httpserver.UploadImage409JSONResponse{
				ConflictJSONResponse: NewConflictError(err),
			}, nil
		default:
			return httpserver.UploadImage500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
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

// DeleteImage Delete image
// (DELETE /api/v1/images/{image_id})
func (h *ImageHandler) DeleteImage(
	ctx context.Context,
	r httpserver.DeleteImageRequestObject,
) (httpserver.DeleteImageResponseObject, error) {
	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return httpserver.DeleteImage401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	if err := h.imageService.Delete(ctx, r.ImageID, user.ID); err != nil {
		switch {
		case errors.Is(err, model.ErrImageNotFound):
			return httpserver.DeleteImage404JSONResponse{
				NotFoundJSONResponse: NewNotFoundError(err),
			}, nil
		default:
			return httpserver.DeleteImage500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.DeleteImage204Response{}, nil
}
