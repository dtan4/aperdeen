package backend

import (
	"github.com/dtan4/aperdeen/model"
)

// Backend represents the interface of real API Gateway
type Backend interface {
	ListEndpoints(apiName string) ([]*model.Endpoint, error)
}
