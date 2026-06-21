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
	return nil, nil
}

// Upload image file
// (POST /api/v1/images)
func (h *ImageHandler) UploadImage(
	ctx context.Context,
	r openapi.UploadImageRequestObject,
) (openapi.UploadImageResponseObject, error) {
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
			BadRequestJSONResponse: openapi.BadRequestJSONResponse{
				Message: "file is required",
			},
		}, nil
	}

	image, err := h.imageService.Create(ctx, files[0])
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnsupportedImageFormat):
			return openapi.UploadImage400JSONResponse{
				BadRequestJSONResponse: openapi.BadRequestJSONResponse{
					Message: err.Error(),
				},
			}, nil
		case errors.Is(err, model.ErrImageAlreadyExists):
			return openapi.UploadImage400JSONResponse{
				BadRequestJSONResponse: openapi.BadRequestJSONResponse{
					Message: err.Error(),
				},
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
// (DELETE /api/v1/images/{id})
func (h *ImageHandler) DeleteImageByID(
	ctx context.Context,
	r openapi.DeleteImageByIDRequestObject,
) (openapi.DeleteImageByIDResponseObject, error) {
	return nil, nil
}
