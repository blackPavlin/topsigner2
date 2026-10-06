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

// GetGroups Get list of groups
// (GET /api/v1/groups)
func (h *GroupHandler) GetGroups(
	ctx context.Context,
	r httpserver.GetGroupsRequestObject,
) (httpserver.GetGroupsResponseObject, error) {
	session, ok := auth.GetSessionFromContext(ctx)
	if !ok {
		return httpserver.GetGroups401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	list, err := h.groupService.List(ctx, session)
	if err != nil {
		switch {
		default:
			return httpserver.GetGroups500JSONResponse{
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

	return httpserver.GetGroups200JSONResponse{
		Items: items,
	}, nil
}

// GetGroupsConnectURL Get URL to connect group(s)
// (GET /api/v1/groups/connect)
func (h *GroupHandler) GetGroupsConnectURL(
	ctx context.Context,
	r httpserver.GetGroupsConnectURLRequestObject,
) (httpserver.GetGroupsConnectURLResponseObject, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return httpserver.GetGroupsConnectURL401JSONResponse{
			UnauthorizedJSONResponse: NewUnauthorizedError(),
		}, nil
	}

	connectionURL, err := h.groupService.GenerateConnectionURL(ctx, user.ID, r.Params.GroupID)
	if err != nil {
		return httpserver.GetGroupsConnectURL500JSONResponse{
			InternalErrorJSONResponse: NewInternalError(),
		}, nil
	}

	return httpserver.GetGroupsConnectURL200JSONResponse{
		URL: connectionURL,
	}, nil
}

// ConnectGroupCallback Connection group callback
// (GET /api/v1/groups/connect/callback)
func (h *GroupHandler) ConnectGroupCallback(
	ctx context.Context,
	r httpserver.ConnectGroupCallbackRequestObject,
) (httpserver.ConnectGroupCallbackResponseObject, error) {
	if _, err := h.groupService.Connect(ctx, r.Params.Code, r.Params.State); err != nil {
		switch {
		case errors.Is(err, model.ErrGroupStateNotFound):
			return httpserver.ConnectGroupCallback400JSONResponse{
				BadRequestJSONResponse: NewBadRequestError(err),
			}, nil
		default:
			return httpserver.ConnectGroupCallback500JSONResponse{
				InternalErrorJSONResponse: NewInternalError(),
			}, nil
		}
	}

	return httpserver.ConnectGroupCallback201Response{}, nil
}
