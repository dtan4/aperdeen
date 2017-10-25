package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/dtan4/aperdeen/service/aws"
	"github.com/dtan4/aperdeen/service/aws/apigateway"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list APINAME",
	Short: "List API endpoints",
	RunE:  doList,
}

var listHeader = []string{
	"PATH",
	"ENDPOINT",
}

func doList(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("APINAME is required")
	}
	apiName := args[0]

	if err := aws.Initialize(rootOpts.region); err != nil {
		return errors.Wrap(err, "cannot create AWS API clients")
	}

	apis, err := aws.APIGateway.ListAPIs()
	if err != nil {
		return errors.Wrap(err, "cannot retrieve APIs")
	}

	var api *apigateway.API

	for _, a := range apis {
		if a.Name == apiName {
			api = a
			break
		}
	}

	if api == nil {
		return errors.Errorf("api %q not found", apiName)
	}

	endpoints, err := aws.APIGateway.ListEndpoints(api.ID)
	if err != nil {
		return errors.Wrap(err, "cannot retrieve endpoints")
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return strings.Compare(endpoints[i].Path, endpoints[j].Path) < 0
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, strings.Join(listHeader, "\t"))

	for _, endpoint := range endpoints {
		fmt.Fprintf(
			w,
			"%s\t%s\n",
			strings.Replace(endpoint.Path, "{proxy+}", "*", -1),
			strings.Replace(endpoint.TargetURL, "{proxy}", "*", -1),
		)
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
