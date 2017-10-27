package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dtan4/aperdeen/model"
	"github.com/gorilla/mux"
)

func TestCreateProxyHandler(t *testing.T) {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	r.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
		case "POST":
			w.WriteHeader(http.StatusCreated)
		}
	}).Methods("GET", "POST")
	r.HandleFunc("/error", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}).Methods("GET")
	r.HandleFunc("/ops/ping", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	// dummy endpoint!
	r.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}).Methods("GET")

	origin := httptest.NewServer(r)
	defer origin.Close()

	endpoints := map[string]*model.Endpoint{
		"/foo/ping": &model.Endpoint{
			Path:      "/foo/ping",
			TargetURL: fmt.Sprintf("%s/ops/ping", origin.URL),
		},
		"/foo/*": &model.Endpoint{
			Path:      "/foo/*",
			TargetURL: fmt.Sprintf("%s/*", origin.URL),
		},
		"/invalid/*": &model.Endpoint{
			Path:      "/invalid/*",
			TargetURL: "thisisainvalidurl",
		},
		"/lambda/*": &model.Endpoint{
			Path:      "/lambda/*",
			TargetURL: "arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/arn:aws:lambda:ap-northeast-1:123456789012:function:api-backend/invocations",
		},
	}

	handler, err := CreateProxyHandler(endpoints)
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}

	s := httptest.NewServer(handler)
	defer s.Close()

	ru, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		method string
		path   string
		want   int
	}{
		// nonexist endpoints
		{
			method: "GET",
			path:   "/bar/",
			want:   http.StatusNotFound,
		},

		{
			method: "GET",
			path:   "/baz/",
			want:   http.StatusNotFound,
		},
		// endpoint exists
		{
			method: "GET",
			path:   "/foo",
			want:   http.StatusOK,
		},
		{
			method: "GET",
			path:   "/foo/",
			want:   http.StatusOK,
		},
		{
			method: "GET",
			path:   "/foo/bar",
			want:   http.StatusOK,
		},
		{
			method: "POST",
			path:   "/foo/bar",
			want:   http.StatusCreated,
		},
		{
			method: "GET",
			path:   "/foo/error",
			want:   http.StatusInternalServerError,
		},
		{
			method: "GET",
			path:   "/foo/ping",
			want:   http.StatusOK,
		},
	}

	c := &http.Client{}

	for _, tc := range testcases {
		ru.Path = tc.path
		req := &http.Request{
			Method: tc.method,
			URL:    ru,
		}

		resp, err := c.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != tc.want {
			t.Errorf("method: %s, path: %s, want: %d, got: %d", tc.method, tc.path, tc.want, resp.StatusCode)
		}
	}
}
