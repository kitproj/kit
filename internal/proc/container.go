package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/strslice"

	"github.com/docker/docker/pkg/stdcopy"

	"github.com/alexec/kit/internal/types"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type container struct {
	types.PodSpec
	types.Task
}

func (c *container) Run(ctx context.Context, stdout, stderr io.Writer) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	if err := c.remove(ctx, cli); err != nil {
		return err
	}

	dockerfile := filepath.Join(c.Image, "Dockerfile")
	if _, err := os.Stat(dockerfile); err == nil {
		r, err := archive.TarWithOptions(filepath.Dir(dockerfile), &archive.TarOptions{})
		if err != nil {
			return err
		}
		defer r.Close()
		resp, err := cli.ImageBuild(ctx, r, dockertypes.ImageBuildOptions{Dockerfile: filepath.Base(dockerfile), Tags: []string{c.Name}})
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err = io.Copy(stdout, resp.Body); err != nil {
			return err
		}
	} else if c.ImagePullPolicy != "Never" {
		r, err := cli.ImagePull(ctx, c.Image, dockertypes.ImagePullOptions{})
		if err != nil {
			return err
		}
		if _, err = io.Copy(stdout, r); err != nil {
			return err
		}
		if err = r.Close(); err != nil {
			return err
		}
	}

	portSet, portBindings, err := c.createPorts()
	if err != nil {
		return err
	}
	binds, err := c.createBinds()
	if err != nil {
		return err
	}
	image := c.Image
	if _, err := os.Stat(filepath.Join(c.Image, "Dockerfile")); err == nil {
		image = c.Name
	}
	created, err := cli.ContainerCreate(ctx, &dockercontainer.Config{
		Hostname:     c.Name,
		ExposedPorts: portSet,
		Tty:          c.TTY,
		Env:          c.Env.Environ(),
		Cmd:          strslice.StrSlice(c.Args),
		Image:        image,
		WorkingDir:   c.WorkingDir,
		// TODO support entrypoint
		Entrypoint: strslice.StrSlice(c.Command),
		Labels:     map[string]string{"name": c.Name},
	}, &dockercontainer.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
	}, &network.NetworkingConfig{}, &v1.Platform{}, c.Name)
	if err != nil {
		return err
	}
	if err = cli.ContainerStart(ctx, created.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		c.remove(context.Background(), cli)
	}()
	logs, err := cli.ContainerLogs(ctx, c.Name, dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	defer logs.Close()
	if _, err = stdcopy.StdCopy(stdout, stderr, logs); err != nil {
		return err
	}
	inspect, err := cli.ContainerInspect(ctx, created.ID)
	if err != nil {
		return err
	}
	if inspect.State.ExitCode > 0 {
		return fmt.Errorf("exit code %d", inspect.State.ExitCode)
	}
	return nil
}

func (c *container) createPorts() (nat.PortSet, map[nat.Port][]nat.PortBinding, error) {
	portSet := nat.PortSet{}
	portBindings := map[nat.Port][]nat.PortBinding{}
	for _, p := range c.Ports {
		port, err := nat.NewPort("tcp", fmt.Sprint(p.ContainerPort))
		if err != nil {
			return nil, nil, err
		}
		portSet[port] = struct{}{}
		hostPort := p.GetHostPort()
		portBindings[port] = []nat.PortBinding{{
			HostPort: fmt.Sprint(hostPort),
		}}
	}
	return portSet, portBindings, nil
}

func (c *container) createBinds() ([]string, error) {
	var binds []string
	for _, mount := range c.VolumeMounts {
		for _, volume := range c.PodSpec.Volumes {
			if volume.Name == mount.Name {
				abs, err := filepath.Abs(volume.HostPath.Path)
				if err != nil {
					return nil, err
				}
				binds = append(binds, fmt.Sprintf("%s:%s", abs, mount.MountPath))
			}
		}
	}
	return binds, nil
}

func (c *container) remove(ctx context.Context, cli *client.Client) error {
	list, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	grace := c.PodSpec.GetTerminationGracePeriod()
	for _, existing := range list {
		if existing.Labels["name"] == c.Name {
			_ = cli.ContainerStop(ctx, existing.ID, &grace)
			if err := cli.ContainerRemove(ctx, existing.ID, dockertypes.ContainerRemoveOptions{Force: true}); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ Interface = &container{}
