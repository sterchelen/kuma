package cmd

import (
	"bytes"
	"context"
	"github.com/Kong/konvoy/components/konvoy-control-plane/api/mesh/v1alpha1"
	"github.com/Kong/konvoy/components/konvoy-control-plane/pkg/core/resources/apis/mesh"
	"github.com/Kong/konvoy/components/konvoy-control-plane/pkg/core/resources/store"
	"github.com/Kong/konvoy/components/konvoy-control-plane/pkg/util/proto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"text/template"
)

type createDataplaneContext struct {
	*rootContext

	args struct {
		Id string
	}
}

func newCreateDataplanesCmd(ctx *createDataplaneContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataplane",
		Short: "Create dataplane",
		Long:  `Create a new Dataplane`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ctx.args.Id == "" {
				return errors.New("Id must not be empty")
			}
			resource, err := parseDataplane(args[0], ctx)
			if err != nil {
				return err
			}
			if err := storeDataplane(ctx, resource); err != nil {
				return err
			}
			return nil
		},
	}
	// flags
	cmd.PersistentFlags().StringVar(&ctx.args.Id, "Id", "", "ID of the dataplane")
	return cmd
}

func parseDataplane(dataplaneTmpl string, ctx *createDataplaneContext) (mesh.DataplaneResource, error) {
	tmpl, err := template.New("").Parse(dataplaneTmpl)
	if err != nil {
		return mesh.DataplaneResource{}, errors.Wrap(err, "could not parse a provided template")
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, ctx.args); err != nil {
		return mesh.DataplaneResource{}, errors.Wrap(err, "could not parse ")
	}
	dp := v1alpha1.Dataplane{}
	if err := proto.FromYAML(buffer.Bytes(), &dp); err != nil {
		return mesh.DataplaneResource{}, errors.Wrap(err, "could not convert yaml to Dataplane")
	}
	// todo also retrieve the meta from yaml
	resource := mesh.DataplaneResource{
		Spec: dp,
	}
	return resource, nil
}

func storeDataplane(ctx *createDataplaneContext, res mesh.DataplaneResource) error {
	controlPlane, err := ctx.CurrentControlPlane()
	if err != nil {
		return err
	}
	rs, err := ctx.NewResourceStore(controlPlane)
	if err != nil {
		return err
	}
	if err := rs.Create(context.Background(), &res, store.CreateByKey(res.Meta.GetNamespace(), res.Meta.GetName(), res.Meta.GetMesh())); err != nil {
		return err
	}
	return nil
}
