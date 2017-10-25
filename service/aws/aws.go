package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	apigatewayapi "github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/pkg/errors"
)

var (
	// APIGateway represents API Gateway Client
	APIGateway *apigateway.Client
)

// Initialize initializes AWS API clients
func Initialize(region string) error {
	var (
		sess *session.Session
		err  error
	)

	if region == "" {
		sess, err = session.NewSession()
		if err != nil {
			return errors.Wrap(err, "Failed to create new AWS session.")
		}
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			return errors.Wrap(err, "Failed to create new AWS session.")
		}
	}

	APIGateway = apigateway.NewClient(apigatewayapi.New(sess))

	return nil
}
