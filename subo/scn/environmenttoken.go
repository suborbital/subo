package scn

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// CreateEnvironmentTokenRequest is a request to create an environment token
type CreateEnvironmentTokenRequest struct {
	Verifier *RequestVerifier
	Env      string `json:"env"`
}

// CreateEnvironmentTokenResponse is a response to a create token request
type CreateEnvironmentTokenResponse struct {
	Token string `json:"token"`
}

// CreateEnvironmentToken creates an environment token
func (s *SCN) CreateEnvironmentToken(verifier *RequestVerifier) (*CreateEnvironmentTokenResponse, error) {
	uri := "/auth/v1/token"
	URL, err := url.Parse(s.endpoint + uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to url.Parse")
	}

	req := &CreateEnvironmentTokenRequest{
		Verifier: verifier,
		Env:      "",
	}

	resp := &CreateEnvironmentTokenResponse{}
	if err := s.Do(http.MethodPost, URL, req, resp); err != nil {
		return nil, errors.Wrap(err, "failed to Do")
	}

	return resp, nil
}
