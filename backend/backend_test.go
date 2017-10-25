package backend

import (
	"testing"
)

func TestConvertEndpointsToYAML(t *testing.T) {
	apiName := "example"
	endpoints := []*Endpoint{
		&Endpoint{
			Path:      "/bar/*",
			TargetURL: "arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations",
		},
		&Endpoint{
			Path:      "/baz/*",
			TargetURL: "",
		},
		&Endpoint{
			Path:      "/foo/*",
			TargetURL: "https://example.com/foo/*",
		},
	}

	got, err := ConvertEndpointsToYAML(apiName, endpoints)
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	want := `name: example
endpoints:
  /bar/*:
    url: arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations
  /baz/*:
    url: ""
  /foo/*:
    url: https://example.com/foo/*
`
	if got != want {
		t.Errorf("want:\n%s\ngot:\n%s", want, got)
	}
}
