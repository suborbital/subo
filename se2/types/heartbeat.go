package types

// HeartbeatRequest is a request to send heartbeat data.
type HeartbeatRequest struct {
	Version string       `json:"version"`
	Modules *ModuleStats `json:"runnables"`
}

// ModuleStats are stats about modules.
type ModuleStats struct {
	TotalCount int `json:"totalCount"`
	IdentCount int `json:"identCount"`
}
