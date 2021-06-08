package scn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

const (
	scnEndpointEnvKey = "SUBO_SCN_ENDPOINT"
)

var defaultEndpoint = "https://api.suborbital.network"

// SCN is an SCN client
type SCN struct {
	endpoint string
}

func New() *SCN {
	endpoint := defaultEndpoint

	if envEndpoint, exists := os.LookupEnv(scnEndpointEnvKey); exists {
		endpoint = envEndpoint
	}

	s := &SCN{
		endpoint: endpoint,
	}

	return s
}

func (s *SCN) Do(method string, URL *url.URL, body, result interface{}) error {
	var buffer io.Reader

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "failed to Marshal")
		}

		buffer = bytes.NewBuffer(bodyJSON)
	}

	request, err := http.NewRequest(method, URL.String(), buffer)
	if err != nil {
		return errors.Wrap(err, "failed to NewRequest")
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to Do request")
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to Do request, received status code %d", resp.StatusCode)
	}

	if result != nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to ReadAll body")
		}

		if err := json.Unmarshal(body, result); err != nil {
			return errors.Wrap(err, "failed to Unmarshal body")
		}
	}

	return nil
}
