package httptransport_test

import (
	"net/http"

	"github.com/bboykiv/topsigner/gen/httpserver"
)

func (s *IntegrationSuite) LoginUser_EmptyBody() {

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := s.client.LoginUserWithResponse(s.T().Context(), httpserver.LoginUserJSONRequestBody{})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode())
	s.Require().Equal(message, resp.JSON400)
}

func (s *IntegrationSuite) LoginUser_EmptyEmail() {

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'required' tag",
	}

	resp, err := s.client.LoginUserWithResponse(s.T().Context(), httpserver.LoginUserJSONRequestBody{
		Password: "password",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode())
	s.Require().Equal(message, resp.JSON400)
}

func (s *IntegrationSuite) LoginUser_InvalidEmail() {

	message := &httpserver.BadRequest{
		Message: "field validation for 'email' failed on the 'email' tag",
	}

	resp, err := s.client.LoginUserWithResponse(s.T().Context(), httpserver.LoginUserJSONRequestBody{
		Email:    "invalid-email",
		Password: "password",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode())
	s.Require().Equal(message, resp.JSON400)
}

func (s *IntegrationSuite) LoginUser_EmptyPassword() {
	message := &httpserver.BadRequest{
		Message: "field validation for 'password' failed on the 'required' tag",
	}

	resp, err := s.client.LoginUserWithResponse(s.T().Context(), httpserver.LoginUserJSONRequestBody{
		Email: "test@email.com",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode())
	s.Require().Equal(message, resp.JSON400)
}
