package se2

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/se2/types"
)

// createEmailVerifier creates an emailverifier (used internally by API.Verify).
func (a *API) createEmailVerifier(email string) (*types.EmailVerifier, error) {
	uri := "/auth/v1/verifier"

	req := &types.CreateEmailVerifierRequest{
		Email: email,
	}

	resp := &types.CreateEmailVerifierResponse{}
	if err := a.do(http.MethodPost, uri, req, resp); err != nil {
		return nil, errors.Wrap(err, "failed to Do")
	}

	return &resp.Verifier, nil
}
