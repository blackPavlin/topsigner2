package e2e_test

import (
	"net/http"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func (s *E2ESuite) TestListImages_Unauthorized() {
	resp, err := s.client.ListImagesWithResponse(s.T().Context(), &httpserver.ListImagesParams{})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode())
	s.Require().Equal(&httpserver.Unauthorized{Message: "unauthorized"}, resp.JSON401)
}
