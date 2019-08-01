package cmd

import (
	"github.com/spf13/cobra"
)

func newCreateCmd(pctx *rootContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Konvoy resources",
		Long:  `Create Konvoy resources.`,
	}
	// sub-commands
	cmd.AddCommand(newCreateDataplanesCmd(&createDataplaneContext{pctx, struct{ Id string }{}}))
	return cmd
}
