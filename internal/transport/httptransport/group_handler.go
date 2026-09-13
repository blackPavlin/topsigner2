package httptransport

import (
	"context"

	"github.com/bboykiv/topsigner/gen/httpserver"
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
			ID:         item.ID,
			Name:       item.Name,
			ScreenName: item.ScreenName,
		})
	}

	return httpserver.GetGroups200JSONResponse{
		Items: items,
	}, nil
}
