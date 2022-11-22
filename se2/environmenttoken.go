package se2

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/se2/types"
)

// CreateEnvironmentToken creates an environment token.
func (a *VerifiedAPI) CreateEnvironmentToken() (*types.CreateEnvironmentTokenResponse, error) {
	uri := "/auth/v1/token"

	req := &types.CreateEnvironmentTokenRequest{
		Verifier: a.verifier,
		Env:      "",
	}

	resp := &types.CreateEnvironmentTokenResponse{}
	if err := a.api.do(http.MethodPost, uri, req, resp); err != nil {
		return nil, errors.Wrap(err, "failed to Do")
	}

	return resp, nil
}
