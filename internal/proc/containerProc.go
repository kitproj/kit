package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	corev1 "k8s.io/api/core/v1"
)

type ContainerProc struct {
	corev1.Container
	cli   *client.Client
	image string
}

func (h *ContainerProc) Init(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	h.cli = cli
	return nil
}

func (h *ContainerProc) Build(ctx context.Context, stdout, stderr io.Writer) error {
	cli := h.cli
	image := h.Image
	if _, err := os.Stat(image); err == nil {
		r, err := archive.TarWithOptions(filepath.Dir(image), &archive.TarOptions{})
		if err != nil {
			return err
		}
		defer r.Close()
		resp, err := cli.ImageBuild(ctx, r, types.ImageBuildOptions{Dockerfile: filepath.Base(image), Tags: []string{h.Name}})
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		_, err = io.Copy(stdout, resp.Body)
		if err != nil {
			return err
		}
		h.image = h.Name
	} else if h.ImagePullPolicy != corev1.PullNever {
		r, err := cli.ImagePull(ctx, h.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		_, err = io.Copy(stdout, r)
		if err != nil {
			return err
		}
		err = r.Close()
		if err != nil {
			return err
		}
		h.image = h.Image
	}
	return nil
}

func (h *ContainerProc) Run(ctx context.Context, stdout, stderr io.Writer) error {

	portSet := nat.PortSet{}
	portBindings := map[nat.Port][]nat.PortBinding{}
	for _, p := range h.Ports {
		port, err := nat.NewPort("tcp", fmt.Sprint(p.ContainerPort))
		if err != nil {
			return err
		}
		portSet[port] = struct{}{}
		hostPort := p.HostPort
		if hostPort == 0 {
			hostPort = p.ContainerPort
		}
		portBindings[port] = []nat.PortBinding{{
			HostPort: fmt.Sprint(hostPort),
		}}
	}

	var environ []string
	for _, env := range h.Env {
		environ = append(environ, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	cli := h.cli
	created, err := cli.ContainerCreate(ctx, &container.Config{
		Hostname: h.Name,
		Env:      environ,
		// TODO support entrypoint
		Entrypoint:   h.Command,
		Cmd:          h.Args,
		Image:        h.image,
		WorkingDir:   h.WorkingDir,
		Tty:          h.TTY,
		ExposedPorts: portSet,
		Labels:       map[string]string{"name": h.Name},
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, &network.NetworkingConfig{}, &v1.Platform{}, h.Name)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	logs, err := cli.ContainerLogs(ctx, h.Name, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	defer logs.Close()
	_, err = io.Copy(stdout, logs)
	return err
}

func (h *ContainerProc) Stop(ctx context.Context, grace time.Duration) error {
	cli := h.cli
	list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	for _, existing := range list {
		if existing.Labels["name"] == h.Name {
			err = cli.ContainerRemove(ctx, existing.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var _ Proc = &ContainerProc{}
