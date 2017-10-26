package backend

import (
	"sort"
	"strings"

	"github.com/dtan4/aperdeen/model"
	"github.com/dtan4/aperdeen/service/aws"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/pkg/errors"
)

// AmazonAPIGateway represents Amazon API Gateway
// implements Backend
type AmazonAPIGateway struct {
	client apigateway.Client
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
func (g *AmazonAPIGateway) ListEndpoints(apiName string) ([]*model.Endpoint, error) {
	apis, err := g.client.ListAPIs()
	if err != nil {
		return []*model.Endpoint{}, errors.Wrap(err, "cannot retrieve APIs")
	}

	var api *apigateway.API

	for _, a := range apis {
		if a.Name == apiName {
			api = a
			break
		}
	}

	if api == nil {
		return []*model.Endpoint{}, errors.Errorf("api %q not found", apiName)
	}

	eps, err := g.client.ListEndpoints(api.ID)
	if err != nil {
		return []*model.Endpoint{}, errors.Wrap(err, "cannot retrieve endpoints")
	}

	endpoints := []*model.Endpoint{}

	for _, ep := range eps {
		endpoints = append(endpoints, &model.Endpoint{
			Path:      strings.Replace(ep.Path, "{proxy+}", "*", -1),
			TargetURL: strings.Replace(ep.TargetURL, "{proxy}", "*", -1),
		})
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return strings.Compare(endpoints[i].Path, endpoints[j].Path) < 0
	})

	return endpoints, nil
}
