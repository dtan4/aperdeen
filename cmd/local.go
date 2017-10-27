package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dtan4/aperdeen/model"
	"github.com/dtan4/aperdeen/proxy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultPort = 8080
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Start local API Gateway",
	RunE:  doLocal,
}

var localOpts = struct {
	filename string
	port     int
}{}

func doLocal(cmd *cobra.Command, args []string) error {
	if localOpts.filename == "" {
		return errors.New("filename (-f, --filename) is required")
	}

	body, err := ioutil.ReadFile(localOpts.filename)
	if err != nil {
		return errors.Wrapf(err, "cannot read %s", localOpts.filename)
	}

	api, err := model.APIFromYAML(body)
	if err != nil {
		return errors.Wrapf(err, "cannot parse %s", localOpts.filename)
	}

	handler, err := proxy.CreateProxyHandler(api.Endpoints)
	if err != nil {
		return errors.Wrap(err, "cannot create proxy handler")
	}

	addr := fmt.Sprintf(":%d", localOpts.port)
	fmt.Printf("server started at %s ...\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		return errors.Wrapf(err, "proxy server at %s returns error", addr)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(localCmd)

	localCmd.Flags().StringVarP(&localOpts.filename, "filename", "f", "", "API definition file")
	localCmd.Flags().IntVarP(&localOpts.port, "port", "p", defaultPort, "local proxy's port number")
}
