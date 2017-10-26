package cmd

import (
	"fmt"

	"github.com/dtan4/aperdeen/backend"
	"github.com/dtan4/aperdeen/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export APINAME",
	Short: "Export API endpoints to YAML",
	RunE:  doExport,
}

func doExport(cmd *cobra.Command, args []string) error {
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

	api := model.BuildAPIWithEndpoints(apiName, endpoints)

	yaml, err := api.ToYAML()
	if err != nil {
		return errors.Wrap(err, "cannot generate API YAML")
	}

	fmt.Printf(yaml)

	return nil
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
