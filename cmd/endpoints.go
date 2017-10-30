package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/aperdeen/backend"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// endpointsCmd represents the endpoints command
var endpointsCmd = &cobra.Command{
	Use:   "endpoints APINAME",
	Short: "List API endpoints",
	RunE:  doEndpoints,
}

var endpointsHeader = []string{
	"PATH",
	"ENDPOINT",
}

func doEndpoints(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("APINAME is required")
	}
	apiName := args[0]

	var be backend.Backend

	be, err := backend.NewAmazonAPIGateway(rootOpts.region)
	if err != nil {
		return errors.Wrap(err, "cannot create AWS API clients")
	}

	endpoints, err := be.ListEndpoints(apiName)
	if err != nil {
		return errors.Wrap(err, "cannot retrieve endpoints")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, strings.Join(endpointsHeader, "\t"))

	for _, endpoint := range endpoints {
		fmt.Fprintf(
			w,
			"%s\t%s\n",
			endpoint.Path,
			endpoint.TargetURL,
		)
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(endpointsCmd)
}
