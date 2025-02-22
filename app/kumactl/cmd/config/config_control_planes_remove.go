package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	kumactl_cmd "github.com/Kong/kuma/app/kumactl/pkg/cmd"
)

func newConfigControlPlanesRemoveCmd(pctx *kumactl_cmd.RootContext) *cobra.Command {
	args := struct {
		name string
	}{}
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a Control Plane",
		Long:  `Remove a Control Plane.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg := pctx.Config()
			if !cfg.RemoveControlPlane(args.name) {
				return errors.Errorf("there is no Control Plane with name %q", args.name)
			}
			if err := pctx.SaveConfig(); err != nil {
				return err
			}
			cmd.Printf("removed Control Plane %q\n", args.name)
			if ctx := cfg.GetCurrent(); ctx != nil {
				cmd.Printf("switched active Control Plane to %q\n", ctx.ControlPlane)
			} else {
				cmd.Printf("there is no active Control Plane left. Use `kumactl config control-planes add` to add a Control Plane and make it active\n")
			}
			return nil
		},
	}
	// flags
	cmd.Flags().StringVar(&args.name, "name", "", "reference name for the Control Plane (required)")
	cmd.MarkFlagRequired("name")
	return cmd
}
