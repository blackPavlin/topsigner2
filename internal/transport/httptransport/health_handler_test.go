package httptransport_test

import (
	"net/http"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func (s *IntegrationSuite) TestHealthCheck_Success() {
	resp, err := s.client.HealthCheckWithResponse(s.T().Context())
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())
	s.Require().Equal(&httpserver.HealthStatus{Status: httpserver.OK}, resp.JSON200)
}
