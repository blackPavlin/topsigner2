package httptransport

import (
	"context"
	"errors"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/group"
)

// todo: добавить проверку прав доступа (для доступа к списку групп нужен тип авторизации VK_OAUTH)

type GroupHandler struct {
	groupService *group.Service
}

func NewGroupHandler(groupService *group.Service) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// ListGroups Get list of groups
// (GET /api/v1/groups)
func (h *GroupHandler) ListGroups(
	ctx context.Context,
	r httpserver.ListGroupsRequestObject,
) (httpserver.ListGroupsResponseObject, error) {
	session, ok := auth.GetSessionFromContext(ctx)
	if !ok {
		return httpserver.ListGroups401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	list, err := h.groupService.List(ctx, session)
	if err != nil {
		switch {
		default:
			return httpserver.ListGroups500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	items := make([]httpserver.Group, 0, len(list))

	for _, item := range list {
		items = append(items, httpserver.Group{
			ID:          item.ID,
			Name:        item.Name,
			ScreenName:  item.ScreenName,
			IsConnected: item.IsConnected,
		})
	}

	return httpserver.ListGroups200JSONResponse{
		Items: items,
	}, nil
}

// GetGroupConnectionURL Get URL to connect group
// (GET /api/v1/groups/connect)
func (h *GroupHandler) GetGroupConnectionURL(
	ctx context.Context,
	r httpserver.GetGroupConnectionURLRequestObject,
) (httpserver.GetGroupConnectionURLResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.GetGroupConnectionURL401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	connectionURL, err := h.groupService.GenerateConnectionURL(ctx, user.ID, r.Params.GroupID)
	if err != nil {
		return httpserver.GetGroupConnectionURL500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
		}, nil
	}

	return httpserver.GetGroupConnectionURL200JSONResponse{
		URL: connectionURL,
	}, nil
}

// HandleGroupCallback Exchange group authorization code
// (GET /api/v1/groups/callback)
func (h *GroupHandler) HandleGroupCallback(
	ctx context.Context,
	r httpserver.HandleGroupCallbackRequestObject,
) (httpserver.HandleGroupCallbackResponseObject, error) {
	// todo: валидация входных параметров

	if _, err := h.groupService.Connect(ctx, *r.Params.Code, *r.Params.State); err != nil {
		switch {
		case errors.Is(err, model.ErrGroupStateNotFound):
			return httpserver.HandleGroupCallback400JSONResponse{
				BadRequestJSONResponse: NewBadRequestError(err),
			}, nil
		default:
			return httpserver.HandleGroupCallback500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.HandleGroupCallback200JSONResponse{}, nil
}
