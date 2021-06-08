package scn

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// EmailVerifier is an email verification record
type EmailVerifier struct {
	ID        int64      `json:"-" db:"id"`
	UUID      string     `json:"uuid" db:"uuid"`
	UserUUID  string     `json:"userUuid" db:"user_uuid"`
	Code      string     `json:"-" db:"code"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	State     string     `json:"state" db:"state"`
}

// RequestVerifier is a verifier used in an HTTP request
type RequestVerifier struct {
	UUID string `json:"uuid"`
	Code string `json:"code"`
}

// CreateEmailVerifierRequest is a request for an email verifier
type CreateEmailVerifierRequest struct {
	Email string `json:"email"`
}

// CreateEmailVerifierResponse is a response to a CreateEmailVerifierRequest
type CreateEmailVerifierResponse struct {
	Verifier EmailVerifier `json:"verifier"`
}

// CreateEmailVerifier creates an emailverifier
func (s *SCN) CreateEmailVerifier(email string) (*EmailVerifier, error) {
	uri := "/auth/v1/verifier"

	URL, err := url.Parse(s.endpoint + uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to url.Parse")
	}

	req := &CreateEmailVerifierRequest{
		Email: email,
	}

	resp := &CreateEmailVerifierResponse{}
	if err := s.Do(http.MethodPost, URL, req, resp); err != nil {
		return nil, errors.Wrap(err, "failed to Do")
	}

	return &resp.Verifier, nil
}
