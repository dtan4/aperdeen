package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	apigatewayapi "github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/pkg/errors"
)

// CreateAPIGatewayClient creates API Gateway client with valid session
func CreateAPIGatewayClient(region string) (*apigateway.Client, error) {
	sess, err := newSession(region)
	if err != nil {
		return nil, errors.Wrap(err, "cannot make new AWS session")
	}

	return apigateway.NewClient(apigatewayapi.New(sess)), nil
}

func newSession(region string) (*session.Session, error) {
	var (
		sess *session.Session
		err  error
	)

	if region == "" {
		sess, err = session.NewSession()
		if err != nil {
			return nil, errors.Wrap(err, "cannot make new AWS session.")
		}
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			return nil, errors.Wrap(err, "cannot make new AWS session.")
		}
	}

	return sess, nil
}
