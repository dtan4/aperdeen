package apigateway

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/apigateway/apigatewayiface"
	"github.com/pkg/errors"
)

// API represents the wrapped form of RestApi
type API struct {
	ID          string
	Name        string
	Description string
	CreatedDate time.Time
}

// Client represents the wrapper of Amazon API Gateway API client
type Client struct {
	api apigatewayiface.APIGatewayAPI
}

// NewClient creates new Client object
func NewClient(api apigatewayiface.APIGatewayAPI) *Client {
	return &Client{
		api: api,
	}
}

// ListAPIs returns the list of registered APIs
func (c *Client) ListAPIs() ([]*API, error) {
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
