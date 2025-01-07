package main

import (
	"context"
	"fmt"
	"io"
	"log"
	url "net/url"
	"os"
	"path/filepath"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kitproj/kit/internal/types"
	"sigs.k8s.io/yaml"
)

type Example struct {
	Name          string         `json:"name"`
	Title         string         `json:"title,omitempty"`
	Uri           string         `json:"uri,omitempty"`
	Description   string         `json:"description,omitempty"`
	Documentation string         `json:"documentation,omitempty"`
	Maintainer    string         `json:"maintainer,omitempty"`
	Pod           types.Workflow `json:"workflow"`
	Licences      string         `json:"licences,omitempty"`
}

func updateExamples(ctx context.Context) error {
	log.Println("updating examples")
	data, err := os.ReadFile("docs/examples/examples.yaml")
	if err != nil {
		return err
	}
	var examples []Example
	if err := yaml.Unmarshal(data, &examples); err != nil {
		return err
	}
	for _, example := range examples {
		if err := updateExample(ctx, &example); err != nil {
			return err
		}
		if err := writeExampleMarkdown(example.Name, example); err != nil {
			return err
		}
		if err := writeExampleYAML(example.Name, example); err != nil {
			return err
		}
	}

	if err := createExamplesReadme(err, examples); err != nil {
		return err
	}
	return nil
}

func createExamplesReadme(err error, examples []Example) error {
	out, err := os.Create("docs/examples/README.md")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString("# Examples\n\n")
	if err != nil {
		return err
	}
	for _, example := range examples {
		description := example.Description
		if description == "" {
			description = example.Documentation
		}
		_, err = out.WriteString(fmt.Sprintf(" * [%s](%s) %s\n", example.Title, example.Name+".md", description))
		if err != nil {
			return err
		}
	}
	return nil
}

func updateExample(ctx context.Context, example *Example) error {
	log.Printf("updating %s", example.Name)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	image := ""
	for _, task := range example.Pod.Tasks {
		image = task.Image
	}
	if err := pullImage(ctx, cli, image); err != nil {
		return err
	}

	inspection, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to inspect %s: %w", image, err)
	}

	example.Title = strings.Title(example.Name)

	// https://github.com/opencontainers/image-spec/blob/main/annotations.md
	for k, v := range inspection.Config.Labels {
		switch k {
		case "org.opencontainers.image.title", "org.label-schema.name", "title", "name":
			example.Title = v
		case "org.opencontainers.image.description", "org.label-schema.description", "description":
			example.Description = v
		case "org.opencontainers.image.documentation", "org.label-schema.usage", "documentation", "usage":
			example.Documentation = v
		case "org.opencontainers.image.url", "org.label-schema.url", "url":
			_, err := url.Parse(v)
			if err == nil {
				example.Uri = v
			}
		case "org.opencontainers.image.licenses", "licence", "license":
			example.Licences = v
		case "org.opencontainers.image.authors", "org.label-schema.vendor", "maintainer", "MAINTAINER", "vendor":
			example.Maintainer = v
		}
	}

	for port := range inspection.Config.ExposedPorts {
		port := port.Int()
		hostPort := port
		if port < 1024 {
			hostPort = 8000 + port
		}

		for _, task := range example.Pod.Tasks {
			task.Ports = append(task.Ports, types.Port{ContainerPort: uint16(port), HostPort: uint16(hostPort)})
		}
	}

	for volume := range inspection.Config.Volumes {
		n := example.Name + "." + filepath.Base(volume)

		for _, task := range example.Pod.Tasks {
			task.VolumeMounts = append(task.VolumeMounts, types.VolumeMount{Name: n, MountPath: volume})
		}
		example.Pod.Volumes = append(example.Pod.Volumes, types.Volume{
			Name:     n,
			HostPath: types.HostPath{Path: filepath.Join("volumes", example.Name, filepath.Base(volume))}})
	}

	return nil
}

func pullImage(ctx context.Context, cli *client.Client, image string) error {
	if resp, err := cli.ImagePull(ctx, image, dockertypes.ImagePullOptions{}); err != nil {
		return fmt.Errorf("failed to pull %s: %w", image, err)
	} else {
		defer resp.Close()
		_, err := io.Copy(os.Stdout, resp)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", image, err)
		}
	}
	return nil
}

func writeExampleYAML(name string, example Example) error {
	out, err := os.Create("docs/examples/" + name + ".yaml")
	if err != nil {
		return err
	}
	defer out.Close()
	data, err := yaml.Marshal(example.Pod)
	if err != nil {
		return err
	}
	if _, err = out.Write(data); err != nil {
		return err
	}
	return nil
}

func writeExampleMarkdown(name string, example Example) error {
	out, err := os.Create("docs/examples/" + name + ".md")
	if err != nil {
		return err
	}
	defer out.Close()

	_, _ = fmt.Fprintf(out, "# %s\n\n", example.Title)
	if example.Uri != "" {
		_, _ = fmt.Fprintf(out, "[Help](%s)\n\n", example.Uri)
	}
	if example.Description != "" {
		_, _ = fmt.Fprintf(out, "%s\n\n", example.Description)
	}
	if example.Documentation != "" {
		_, _ = fmt.Fprintf(out, "%s\n\n", example.Documentation)
	}
	if example.Maintainer != "" {
		_, _ = fmt.Fprintf(out, "> Maintainer: %s\n\n", example.Maintainer)
	}
	data, _ := yaml.Marshal(example.Pod)
	_, _ = fmt.Fprintf(out, "```yaml\n%s```\n\n", string(data))
	if example.Licences != "" {
		_, _ = fmt.Fprintf(out, "Licence(s): %s\n\n", example.Licences)
	}
	return nil
}
