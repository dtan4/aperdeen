package backend

// Endpoint represents API endpoint
type Endpoint struct {
	Path      string
	TargetURL string
}

// Backend represents the interface of real API Gateway
type Backend interface {
	ListEndpoints(apiName string) ([]*Endpoint, error)
}
