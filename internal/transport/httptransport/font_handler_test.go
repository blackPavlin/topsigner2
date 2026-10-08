package httptransport_test

import (
	"net/http"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func (s *IntegrationSuite) TestListFonts_Unauthorized() {
	resp, err := s.client.ListFontsWithResponse(s.T().Context(), &httpserver.ListFontsParams{})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode())
	s.Require().Equal(&httpserver.Unauthorized{Message: "unauthorized"}, resp.JSON401)
}
