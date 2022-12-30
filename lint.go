package main

import (
	"os"

	"github.com/alexec/kit/internal/types"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func lint() *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "Lint file",
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := os.ReadFile(kitFile)
			if err != nil {
				return err
			}
			pod := &types.Pod{}
			err = yaml.UnmarshalStrict(in, pod)
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(pod)
			if err != nil {
				return err
			}
			err = os.WriteFile(kitFile, data, 0o0644)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
