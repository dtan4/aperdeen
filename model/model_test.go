package model

import (
	"reflect"
	"testing"
)

func TestAPIFromYAML(t *testing.T) {
	body := []byte(`name: acme-api
endpoints:
  "/foo/ping":
    url: "https://foo${variable.prefix}.example.com/ping"
  "/foo/*":
    url: "https://foo${variable.prefix}.example.com/*"
  "/bar/*":
    url: "https://bar.example.com/*"
`)

	got, err := APIFromYAML(body)
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	want := &API{
		Name: "acme-api",
		Endpoints: map[string]*Endpoint{
			"/foo/ping": &Endpoint{
				Path:      "/foo/ping",
				TargetURL: "https://foo${variable.prefix}.example.com/ping",
			},
			"/foo/*": &Endpoint{
				Path:      "/foo/*",
				TargetURL: "https://foo${variable.prefix}.example.com/*",
			},
			"/bar/*": &Endpoint{
				Path:      "/bar/*",
				TargetURL: "https://bar.example.com/*",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nwant: %s\ngot:  %s", want, got)
	}
}

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
