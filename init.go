package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/alexec/kit/internal/types"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

//go:embed kit.yaml
var kitYaml string

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "init",
		Long: "Initialize config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := defaultConfigFile
			pod := &types.Pod{}
			err := yaml.Unmarshal([]byte(kitYaml), pod)
			if err != nil {
				return err
			}

			wd, _ := os.Getwd()
			pod.Metadata.Name = filepath.Base(wd)

			data, err := yaml.Marshal(pod)
			if err != nil {
				return err
			}
			err = os.WriteFile(configFile, data, 0o0644)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
