package backend

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

// Backend represents the interface of real API Gateway
type Backend interface {
	ListEndpoints(apiName string) ([]*Endpoint, error)
}

// ConvertEndpointsToYAML converts the given API spec to YAML
func ConvertEndpointsToYAML(apiName string, endpoints []*Endpoint) (string, error) {
	eps := map[string]*Endpoint{}

	for _, ep := range endpoints {
		eps[ep.Path] = ep
	}

	api := &API{
		Name:      apiName,
		Endpoints: eps,
	}

	d, err := yaml.Marshal(api)
	if err != nil {
		return "", errors.Wrap(err, "cannot convert to YAML")
	}

	return string(d), nil
}
