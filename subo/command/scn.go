package command

import (
	"os"

	"github.com/suborbital/subo/scn"
)

func scnAPI() *scn.API {
	endpoint := scn.DefaultEndpoint
	if envEndpoint, exists := os.LookupEnv(scnEndpointEnvKey); exists {
		endpoint = envEndpoint
	}

	api := scn.New(endpoint)

	return api
}
