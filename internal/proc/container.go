package proc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/kitproj/kit/internal/types"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"k8s.io/utils/strings/slices"
)

type container struct {
	types.PodSpec
	types.Task
}

func (c *container) Run(ctx context.Context, stdout, stderr io.Writer) error {

	data, _ := json.Marshal(c.Task)
	expectedHash := base64.StdEncoding.EncodeToString(sha256.New().Sum(data))

	log.Printf("%s: expected hash: %s", c.Name, expectedHash)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	dockerfile := filepath.Join(c.Image, "Dockerfile")
	id, existingHash, err := c.getContainer(ctx, cli)

	log.Printf("%s: existing hash: %s", c.Name, existingHash)

	// If the container exists and the hash is different, remove it.
	if id != "" && existingHash != expectedHash {
		log.Printf("%s: removing container %s", c.Name, id)
		if err := cli.ContainerRemove(ctx, id, dockertypes.ContainerRemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}
		id = ""
	}

	if err != nil {
		return fmt.Errorf("failed to get container ID: %w", err)
	} else if id != "" {
		log.Printf("%s: container already exists, skipping build/pull\n", c.Name)
	} else if _, err := os.Stat(dockerfile); err == nil {
		log.Printf("%s: creating tar image from %q", c.Name, dockerfile)
		r, err := archive.TarWithOptions(filepath.Dir(dockerfile), &archive.TarOptions{})
		if err != nil {
			return err
		}
		defer r.Close()
		log.Printf("%s: building image from %q", c.Name, dockerfile)
		resp, err := cli.ImageBuild(ctx, r, dockertypes.ImageBuildOptions{Dockerfile: filepath.Base(dockerfile), Tags: []string{c.Name}})
		if err != nil {
			return fmt.Errorf("failed to build image: %w", err)
		}
		defer resp.Body.Close()
		log.Printf("%s: building image from %q (logs)", c.Name, dockerfile)
		if _, err = io.Copy(stdout, resp.Body); err != nil {
			return fmt.Errorf("failed to build image (logs): %w", err)
		}
	} else if c.ImagePullPolicy != "Never" {
		log.Printf("%s: pulling image %q", c.Name, c.Image)
		r, err := cli.ImagePull(ctx, c.Image, dockertypes.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
		if _, err = io.Copy(stdout, r); err != nil {
			return fmt.Errorf("failed to pull image (logs): %w", err)
		}
		if err = r.Close(); err != nil {
			return fmt.Errorf("failed to pull image (close): %w", err)
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

	log.Printf("%s: creating container", c.Name)
	_, err = cli.ContainerCreate(ctx, &dockercontainer.Config{
		Hostname:     c.Name,
		ExposedPorts: portSet,
		Tty:          c.TTY,
		Env:          c.Env.Environ(),
		Cmd:          strslice.StrSlice(c.Args),
		Image:        image,
		User:         c.User,
		WorkingDir:   c.WorkingDir,
		Entrypoint:   strslice.StrSlice(c.Command),
		Labels:       map[string]string{hashLabel: expectedHash},
	}, &dockercontainer.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
	}, &network.NetworkingConfig{}, &v1.Platform{}, c.Name)
	if ignoreConflict(err) != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	id, _, err = c.getContainer(ctx, cli)
	if err != nil {
		return err
	}
	if err = cli.ContainerStart(ctx, id, dockertypes.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	go func() {
		<-ctx.Done()
		log.Printf("%s: context cancelled, stopping container", c.Name)
		if err := c.stop(context.Background()); err != nil {
			log.Printf("%s: failed to stop: %v", c.Name, err)
		}
	}()
	log.Printf("%s: logging container\n", c.Name)
	logs, err := cli.ContainerLogs(ctx, c.Name, dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Since:      time.Now().Format(time.RFC3339),
	})
	log.Printf("%s: logged container %v", c.Name, err)
	if err != nil {
		return fmt.Errorf("failed to log container: %w", err)
	}
	defer logs.Close()
	if _, err = stdcopy.StdCopy(stdout, stderr, logs); err != nil {
		return err
	}
	waitC, errC := cli.ContainerWait(ctx, id, dockercontainer.WaitConditionNotRunning)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case wait := <-waitC:
		if wait.StatusCode != 0 {
			return fmt.Errorf("exit code %d", wait.StatusCode)
		}
		return nil
	case err := <-errC:
		return err
	}
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

func (c *container) Reset(ctx context.Context) error {
	return c.stop(ctx)
}

func (c *container) stop(ctx context.Context) error {
	log.Printf("%s: stopping container\n", c.Name)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()
	id, _, err := c.getContainer(ctx, cli)
	if err != nil {
		return err
	}
	grace := c.PodSpec.GetTerminationGracePeriod()
	log.Printf("%s: stopping container %q\n", c.Name, id)
	err = cli.ContainerStop(ctx, id, &grace)
	log.Printf("%s: stopped container %q: %v\n", c.Name, id, err)
	if ignoreNotExist(err) != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

const hashLabel = "kit.hash"

func (c *container) getContainer(ctx context.Context, cli *client.Client) (string, string, error) {
	log.Printf("%s: listing containers", c.Name)
	list, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
	log.Printf("%s: listed %d containers: %v", c.Name, len(list), err)
	if err != nil {
		return "", "", err
	}
	for _, existing := range list {
		if slices.Contains(existing.Names, "/"+c.Name) {
			id := existing.ID
			log.Printf("%s: found container: %s", c.Name, id)
			return id, existing.Labels[hashLabel], nil
		}
	}
	log.Printf("%s: container not found", c.Name)
	return "", "", nil
}

func ignoreConflict(err error) error {
	if errdefs.IsConflict(err) {
		return nil
	}
	return err
}

func ignoreNotExist(err error) error {
	if errdefs.IsNotFound(err) {
		return nil
	}
	return err

}

var _ Interface = &container{}
