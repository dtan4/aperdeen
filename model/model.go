package model

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

/*
name: acme-api
auth_function:
  prod: <function arn>
  qa: <function arn>
  local: <function arn>
endpoints:
  "/foo/ping":
    auth: false
    url: "https://foo${variable.prefix}.example.com/ping"
  "/foo/*":
    auth: true
    url: "https://foo${variable.prefix}.example.com/*"
  "/bar/*":
    auth: true
    url: "https://bar.example.com/*"
*/

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

// APIFromYAML creates API object from YAML
func APIFromYAML(body []byte) (*API, error) {
	var api API

	if err := yaml.Unmarshal(body, &api); err != nil {
		return nil, errors.Wrap(err, "cannot parse the given YAML as API")
	}

	for k := range api.Endpoints {
		api.Endpoints[k].Path = k
	}

	return &api, nil
}

// ToYAML converts API object to YAML
func (a *API) ToYAML() (string, error) {
	d, err := yaml.Marshal(a)
	if err != nil {
		return "", errors.Wrap(err, "cannot convert to YAML")
	}

	return string(d), nil
}
