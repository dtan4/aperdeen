package backend

import (
	"reflect"
	"testing"
	"time"

	"github.com/dtan4/aperdeen/model"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/golang/mock/gomock"
)

func TestListEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := apigateway.NewMockClient(ctrl)
	client.EXPECT().ListAPIs().Return([]*apigateway.API{
		&apigateway.API{
			ID:          "abcde12345",
			Name:        "example",
			Description: "example APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC(),
		},
		&apigateway.API{
			ID:          "12345abcde",
			Name:        "foobar",
			Description: "foobar APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC(),
		},
	}, nil)
	client.EXPECT().ListEndpoints("abcde12345").Return([]*apigateway.Endpoint{
		&apigateway.Endpoint{
			Path:      "/foo/{proxy+}",
			TargetURL: "https://example.com/foo/{proxy}",
		},
		&apigateway.Endpoint{
			Path:      "/bar/{proxy+}",
			TargetURL: "arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations",
		},
		&apigateway.Endpoint{
			Path:      "/baz/{proxy+}",
			TargetURL: "",
		},
	}, nil)
	gw := &AmazonAPIGateway{
		client: client,
	}

	got, err := gw.ListEndpoints("example")
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	want := []*model.Endpoint{
		&model.Endpoint{
			Path:      "/bar/*",
			TargetURL: "arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations",
		},
		&model.Endpoint{
			Path:      "/baz/*",
			TargetURL: "",
		},
		&model.Endpoint{
			Path:      "/foo/*",
			TargetURL: "https://example.com/foo/*",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestListEndpoints_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := apigateway.NewMockClient(ctrl)
	client.EXPECT().ListAPIs().Return([]*apigateway.API{
		&apigateway.API{
			ID:          "abcde12345",
			Name:        "example",
			Description: "example APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 34, 56, 0, time.UTC).UTC(),
		},
		&apigateway.API{
			ID:          "12345abcde",
			Name:        "foobar",
			Description: "foobar APIs",
			CreatedDate: time.Date(2017, 10, 25, 12, 00, 00, 0, time.UTC).UTC(),
		},
	}, nil)
	gw := &AmazonAPIGateway{
		client: client,
	}

	_, err := gw.ListEndpoints("baz")
	if err == nil {
		t.Errorf("got no error")
		return
	}

	want := "api \"baz\" not found"

	if err.Error() != want {
		t.Errorf("got: %q, want: %q", err.Error(), want)
	}
}
