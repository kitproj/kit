package proc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	types.Spec
	types.Container
	cli *client.Client
}

func (h *container) Init(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	h.cli = cli
	return nil
}

func (h *container) Build(ctx context.Context, stdout, stderr io.Writer) error {

	if build := h.Container.Build; build != nil {
		if build.HasMutex() {
			if _, err := stdout.Write([]byte(fmt.Sprintf("waiting for mutex %q to unlock...\n", build.Mutex))); err != nil {
				return err
			}
			mutex := KeyLock(build.Mutex)
			mutex.Lock()
			defer mutex.Unlock()
			if _, err := stdout.Write([]byte(fmt.Sprintf("locked mutex %q\n", build.Mutex))); err != nil {
				return err
			}
		}
	}
	dockerfile := filepath.Join(h.Image, "Dockerfile")
	if _, err := os.Stat(dockerfile); err == nil {
		r, err := archive.TarWithOptions(filepath.Dir(dockerfile), &archive.TarOptions{})
		if err != nil {
			return err
		}
		defer r.Close()
		resp, err := h.cli.ImageBuild(ctx, r, dockertypes.ImageBuildOptions{Dockerfile: filepath.Base(dockerfile), Tags: []string{h.Image}})
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err = io.Copy(stdout, resp.Body); err != nil {
			return err
		}
	} else if h.ImagePullPolicy != "PullNever" {
		r, err := h.cli.ImagePull(ctx, h.Image, dockertypes.ImagePullOptions{})
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
	return nil
}

func (h *container) Run(ctx context.Context, stdout, stderr io.Writer) error {
	if err := h.remove(ctx); err != nil {
		return err
	}
	portSet, portBindings, err := h.createPorts()
	if err != nil {
		return err
	}
	binds, err := h.createBinds()
	if err != nil {
		return err
	}
	created, err := h.cli.ContainerCreate(ctx, &dockercontainer.Config{
		Hostname:     h.Name,
		ExposedPorts: portSet,
		Tty:          h.TTY,
		Env:          h.Env.Environ(),
		Cmd:          h.Args,
		Image:        h.Image,
		WorkingDir:   h.WorkingDir,
		// TODO support entrypoint
		Entrypoint: h.Command,
		Labels:     map[string]string{"name": h.Name},
	}, &dockercontainer.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
	}, &network.NetworkingConfig{}, &v1.Platform{}, h.Name)
	if err != nil {
		return err
	}
	if err = h.cli.ContainerStart(ctx, created.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		h.remove(context.Background())
	}()
	logs, err := h.cli.ContainerLogs(ctx, h.Name, dockertypes.ContainerLogsOptions{
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
	inspect, err := h.cli.ContainerInspect(ctx, created.ID)
	if err != nil {
		return err
	}
	if inspect.State.ExitCode > 0 {
		return fmt.Errorf("exit code %d", inspect.State.ExitCode)
	}
	return nil
}

func (h *container) createPorts() (nat.PortSet, map[nat.Port][]nat.PortBinding, error) {
	portSet := nat.PortSet{}
	portBindings := map[nat.Port][]nat.PortBinding{}
	for _, p := range h.Ports {
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

func (h *container) createBinds() ([]string, error) {
	var binds []string
	for _, mount := range h.VolumeMounts {
		for _, volume := range h.Spec.Volumes {
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

func (h *container) remove(ctx context.Context) error {
	list, err := h.cli.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	grace := h.Spec.GetTerminationGracePeriod()
	for _, existing := range list {
		if existing.Labels["name"] == h.Name {
			_ = h.cli.ContainerStop(ctx, existing.ID, &grace)
			if err := h.cli.ContainerRemove(ctx, existing.ID, dockertypes.ContainerRemoveOptions{Force: true}); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ Interface = &container{}
