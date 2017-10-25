package apigateway

import (
	"github.com/aws/aws-sdk-go/service/apigateway/apigatewayiface"
)

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
