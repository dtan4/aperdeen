package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/dtan4/aperdeen/model"
	"github.com/gorilla/mux"
)

var availableSchemas = map[string]bool{
	"http":  true,
	"https": true,
}

// CreateProxyHandler creates local API Gateway server
func CreateProxyHandler(endpoints map[string]*model.Endpoint) (http.Handler, error) {
	r := mux.NewRouter().StrictSlash(true)

	for p, ep := range endpoints {
		target, err := url.Parse(ep.TargetURL)
		if err != nil {
			// return nil, errors.Wrapf(err, "cannot parse %q as URL", ep.TargetURL)
			continue
		}

		if _, ok := availableSchemas[target.Scheme]; !ok {
			// return nil, errors.Errorf("invalid URI schema: %s", target.Scheme)
			continue
		}

		d := director(target)

		// add "/foo/{rest:.*}" route from "/foo/*"
		r.Handle(strings.Replace(p, "*", "{rest:.*}", -1), &httputil.ReverseProxy{
			Director: d,
		})

		// add "/foo" route from "/foo/*"
		if strings.HasSuffix(p, "/*") {
			r.Handle(strings.TrimSuffix(p, "/*"), &httputil.ReverseProxy{
				Director: d,
			})
		}
	}

	return r, nil
}

func director(target *url.URL) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = strings.Replace(target.Path, "*", mux.Vars(req)["rest"], -1)

		targetQuery := target.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}
