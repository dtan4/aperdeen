package model

import (
	"testing"
)

func TestToYAML(t *testing.T) {
	api := &API{
		Name: "example",
		Endpoints: map[string]*Endpoint{
			"/bar/*": &Endpoint{
				Path:      "/bar/*",
				TargetURL: "arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations",
			},
			"/baz/*": &Endpoint{
				Path:      "/baz/*",
				TargetURL: "",
			},
			"/foo/*": &Endpoint{
				Path:      "/foo/*",
				TargetURL: "https://example.com/foo/*",
			},
		},
	}

	got, err := api.ToYAML()
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
