package se2

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/suborbital/subo/se2/types"
)

// SendHeartbeat sends a telemetry heartbeat request.
func (e *EnvironmentAPI) SendHeartbeat(heartbeat *types.HeartbeatRequest) error {
	uri := "/telemetry/v1/heartbeat"

	headers := map[string]string{
		tokenRequestHeaderKey: e.token,
	}

	if err := e.api.doWithHeaders(http.MethodPost, uri, headers, heartbeat, nil); err != nil {
		return errors.Wrap(err, "failed to doWithHeaders")
	}

	return nil
}
