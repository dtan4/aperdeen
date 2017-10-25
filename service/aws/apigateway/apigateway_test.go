package apigateway

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
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

func TestListAPIs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewMockAPIGatewayAPI(ctrl)
	api.EXPECT().GetRestApis(&apigateway.GetRestApisInput{}).Return(&apigateway.GetRestApisOutput{
		Items: []*apigateway.RestApi{
			&apigateway.RestApi{
				Id:          aws.String("abcde12345"),
				Name:        aws.String("example"),
				Description: aws.String("example APIs"),
				CreatedDate: aws.Time(time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC()),
			},
			&apigateway.RestApi{
				Id:          aws.String("12345abcde"),
				Name:        aws.String("foobar"),
				Description: aws.String("foobar APIs"),
				CreatedDate: aws.Time(time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC()),
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	got, err := client.ListAPIs()
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	want := []*API{
		&API{
			ID:          "abcde12345",
			Name:        "example",
			Description: "example APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC(),
		},
		&API{
			ID:          "12345abcde",
			Name:        "foobar",
			Description: "foobar APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC(),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

func TestListAPIs_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewMockAPIGatewayAPI(ctrl)
	api.EXPECT().GetRestApis(&apigateway.GetRestApisInput{}).Return(&apigateway.GetRestApisOutput{}, fmt.Errorf("error"))
	client := &Client{
		api: api,
	}

	_, err := client.ListAPIs()
	if err == nil {
		t.Errorf("got no error")
		return
	}

	want := "cannot retrieve registered APIs: error"

	if err.Error() != want {
		t.Errorf("got: %#v, want: %#v", err.Error(), want)
	}
}

func TestListStages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewMockAPIGatewayAPI(ctrl)
	api.EXPECT().GetStages(&apigateway.GetStagesInput{
		RestApiId: aws.String("abcde12345"),
	}).Return(&apigateway.GetStagesOutput{
		Item: []*apigateway.Stage{
			&apigateway.Stage{
				StageName:       aws.String("prod"),
				DeploymentId:    aws.String("abc123"),
				LastUpdatedDate: aws.Time(time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC()),
			},
			&apigateway.Stage{
				StageName:       aws.String("qa"),
				DeploymentId:    aws.String("123abc"),
				LastUpdatedDate: aws.Time(time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC()),
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	got, err := client.ListStages("abcde12345")
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	want := []*Stage{
		&Stage{
			Name:            "prod",
			DeploymentID:    "abc123",
			LastUpdatedDate: time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC(),
		},
		&Stage{
			Name:            "qa",
			DeploymentID:    "123abc",
			LastUpdatedDate: time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC(),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

func TestListStages_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewMockAPIGatewayAPI(ctrl)
	api.EXPECT().GetStages(&apigateway.GetStagesInput{
		RestApiId: aws.String("abcde12345"),
	}).Return(&apigateway.GetStagesOutput{}, fmt.Errorf("error"))
	client := &Client{
		api: api,
	}

	_, err := client.ListStages("abcde12345")
	if err == nil {
		t.Errorf("got no error")
		return
	}

	want := "cannot retrieve API stages: error"

	if err.Error() != want {
		t.Errorf("got: %#v, want: %#v", err.Error(), want)
	}
}
