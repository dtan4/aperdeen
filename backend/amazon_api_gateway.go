package backend

import (
	"sort"
	"strings"

	"github.com/dtan4/aperdeen/service/aws"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/pkg/errors"
)

// AmazonAPIGateway represents Amazon API Gateway
// implements Backend
type AmazonAPIGateway struct {
	client *apigateway.Client
}

// NewAmazonAPIGateway creates new AmazonAPIGateway object
func NewAmazonAPIGateway(region string) (*AmazonAPIGateway, error) {
	client, err := aws.CreateAPIGatewayClient(region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create API Gateway client")
	}

	return &AmazonAPIGateway{
		client: client,
	}, nil
}

// ListEndpoints returns the list of registered API endpoints
func (g *AmazonAPIGateway) ListEndpoints(apiName string) ([]*Endpoint, error) {
	apis, err := g.client.ListAPIs()
	if err != nil {
		return []*Endpoint{}, errors.Wrap(err, "cannot retrieve APIs")
	}

	var api *apigateway.API

	for _, a := range apis {
		if a.Name == apiName {
			api = a
			break
		}
	}

	if api == nil {
		return []*Endpoint{}, errors.Errorf("api %q not found", apiName)
	}

	eps, err := g.client.ListEndpoints(api.ID)
	if err != nil {
		return []*Endpoint{}, errors.Wrap(err, "cannot retrieve endpoints")
	}

	endpoints := []*Endpoint{}

	for _, ep := range eps {
		endpoints = append(endpoints, &Endpoint{
			Path:      ep.Path,
			TargetURL: ep.TargetURL,
		})
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return strings.Compare(endpoints[i].Path, endpoints[j].Path) < 0
	})

	return endpoints, nil
}
