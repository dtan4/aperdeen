package model

// API represents a entire model of API
type API struct {
	Name      string               `yaml:"name"`
	Endpoints map[string]*Endpoint `yaml:"endpoints"`
}

// Endpoint represents API endpoint
type Endpoint struct {
	Path      string `yaml:"-"`
	TargetURL string `yaml:"url"`
}

// BuildAPIWithEndpoints craete new build object with the given endpoints
func BuildAPIWithEndpoints(apiName string, endpoints []*Endpoint) *API {
	eps := map[string]*Endpoint{}

	for _, ep := range endpoints {
		eps[ep.Path] = ep
	}

	return &API{
		Name:      apiName,
		Endpoints: eps,
	}
}
