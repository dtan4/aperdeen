package apigateway

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewMockAPIGatewayAPI(ctrl)

	client := NewClient(api)
	if client.api != api {
		t.Error("invalid client")
	}
}
