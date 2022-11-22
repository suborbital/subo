package command

import (
	"os"

	"github.com/suborbital/subo/se2"
)

func se2API() *se2.API {
	endpoint := se2.DefaultEndpoint
	if envEndpoint, exists := os.LookupEnv(se2EndpointEnvKey); exists {
		endpoint = envEndpoint
	}

	api := se2.New(endpoint)

	return api
}
