//go:generate mockgen -source=apigateway.go -destination=apigateway_mock.go -package=apigateway

package apigateway

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/apigateway/apigatewayiface"
	"github.com/pkg/errors"
)

const (
	httpMethodANY = "ANY"
)

// API represents the wrapped form of RestApi
type API struct {
	ID          string
	Name        string
	Description string
	CreatedDate time.Time
}

// Endpoint represents API endpoint
type Endpoint struct {
	Path      string
	TargetURL string
}

// Stage represents the wrapped form of Stage
type Stage struct {
	Name            string
	DeploymentID    string
	LastUpdatedDate time.Time
}

// Client is the wrapper of Amazon API Gateway API client
type Client interface {
	ListAPIs() ([]*API, error)
	ListEndpoints(apiID string) ([]*Endpoint, error)
	ListStages(apiID string) ([]*Stage, error)
}

// ClientImpl represents the implementation of Client
type ClientImpl struct {
	api apigatewayiface.APIGatewayAPI
}

// NewClientImpl creates new ClientImpl object
func NewClientImpl(api apigatewayiface.APIGatewayAPI) *ClientImpl {
	return &ClientImpl{
		api: api,
	}
}

// ListAPIs returns the list of registered APIs
func (c *ClientImpl) ListAPIs() ([]*API, error) {
	resp, err := c.api.GetRestApis(&apigateway.GetRestApisInput{})
	if err != nil {
		return []*API{}, errors.Wrap(err, "cannot retrieve registered APIs")
	}

	apis := []*API{}

	for _, item := range resp.Items {
		apis = append(apis, &API{
			ID:          aws.StringValue(item.Id),
			Name:        aws.StringValue(item.Name),
			Description: aws.StringValue(item.Description),
			CreatedDate: aws.TimeValue(item.CreatedDate),
		})
	}

	return apis, nil
}

// ListEndpoints returne the endpoints of the given API
func (c *ClientImpl) ListEndpoints(apiID string) ([]*Endpoint, error) {
	resources, err := c.api.GetResources(&apigateway.GetResourcesInput{
		RestApiId: aws.String(apiID),
	})
	if err != nil {
		return []*Endpoint{}, errors.Wrap(err, "cannot retrieve API resources")
	}

	endpoints := []*Endpoint{}

	for _, r := range resources.Items {
		if _, ok := r.ResourceMethods[httpMethodANY]; !ok {
			continue
		}

		integration, err := c.api.GetIntegration(&apigateway.GetIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: r.Id,
			HttpMethod: aws.String(httpMethodANY),
		})
		// skip error to include URL-unassigned endpoints in response
		if err == nil {
			endpoints = append(endpoints, &Endpoint{
				Path:      aws.StringValue(r.Path),
				TargetURL: aws.StringValue(integration.Uri),
			})
		} else {
			endpoints = append(endpoints, &Endpoint{
				Path:      aws.StringValue(r.Path),
				TargetURL: "",
			})
		}
	}

	return endpoints, nil
}

// ListStages returns the list of registered APIs
func (c *ClientImpl) ListStages(apiID string) ([]*Stage, error) {
	resp, err := c.api.GetStages(&apigateway.GetStagesInput{
		RestApiId: aws.String(apiID),
	})
	if err != nil {
		return []*Stage{}, errors.Wrap(err, "cannot retrieve API stages")
	}

	stages := []*Stage{}

	for _, item := range resp.Item {
		stages = append(stages, &Stage{
			Name:            aws.StringValue(item.StageName),
			DeploymentID:    aws.StringValue(item.DeploymentId),
			LastUpdatedDate: aws.TimeValue(item.LastUpdatedDate),
		})
	}

	return stages, nil
}
