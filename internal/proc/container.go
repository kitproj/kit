package proc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/registry"
	"github.com/docker/go-connections/nat"
	"github.com/kitproj/kit/internal/metrics"
	"github.com/kitproj/kit/internal/types"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"k8s.io/utils/strings/slices"
)

type container struct {
	name string
	log  *log.Logger
	spec types.Spec
	types.Task
	containerID    string
	profFSSnapshot *metrics.ProcFSSnapshot
}

func (c *container) Run(ctx context.Context, stdout, stderr io.Writer) error {

	log := c.log
	data, _ := json.Marshal(c.Task)
	expectedHash := fmt.Sprintf("%x", adler32.Checksum(data))

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	dockerfile := filepath.Join(c.Image, "Dockerfile")
	id, existingHash, err := c.getContainer(ctx, cli)

	// If the container exists and the hash is different, remove it.
	if id != "" && existingHash != expectedHash {
		log.Println("removing container")
		if err := cli.ContainerRemove(ctx, id, dockertypes.ContainerRemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}
		id = ""
	}

	environ, err := types.Environ(c.spec, c.Task)
	if err != nil {
		return fmt.Errorf("error getting spec environ: %w", err)
	}

	if id != "" {
		log.Printf("container already exists, skipping build/pull\n")
	} else if _, err := os.Stat(dockerfile); err == nil {
		log.Printf("creating tar image from %q", dockerfile)
		r, err := archive.TarWithOptions(filepath.Dir(dockerfile), &archive.TarOptions{})
		if err != nil {
			return fmt.Errorf("failed to create tar: %w", err)
		}
		defer r.Close()
		log.Printf("building image from %q", dockerfile)
		resp, err := cli.ImageBuild(ctx, r, dockertypes.ImageBuildOptions{Dockerfile: filepath.Base(dockerfile), Tags: []string{c.name}})
		if err != nil {
			return fmt.Errorf("failed to build image: %w", err)
		}
		defer resp.Body.Close()
		log.Printf("building image from %q (logs)", dockerfile)
		if _, err = io.Copy(stdout, resp.Body); err != nil {
			return fmt.Errorf("failed to build image (logs): %w", err)
		}
	} else if c.ImagePullPolicy != "Never" {
		log.Printf("pulling image %q", c.Image)

		ref, err := reference.ParseNormalizedNamed(c.Image)
		if err != nil {
			return fmt.Errorf("unable to parse image: %w", err)
		}
		repoInfo, err := registry.ParseRepositoryInfo(ref)
		if err != nil {
			return fmt.Errorf("unable to parse repository info: %w", err)
		}

		var server string
		if repoInfo.Index.Official {
			info, err := cli.Info(ctx)
			if err != nil || info.IndexServerAddress == "" {
				server = registry.IndexServer
			} else {
				server = info.IndexServerAddress
			}
		} else {
			server = repoInfo.Index.Name
		}
		errBuf := &bytes.Buffer{}
		cf := config.LoadDefaultConfigFile(errBuf)
		if errBuf.Len() > 0 {
			return fmt.Errorf("unable to load docker config: %s", errBuf.String())
		}
		authConfig, err := cf.GetAuthConfig(server)
		if err != nil {
			return fmt.Errorf("failed to get auth config: %w", err)
		}
		buf, err := json.Marshal(authConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal auth config: %w", err)
		}
		encodedAuth := base64.URLEncoding.EncodeToString(buf)

		r, err := cli.ImagePull(ctx, c.Image, dockertypes.ImagePullOptions{
			RegistryAuth: encodedAuth,
		})
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
		return fmt.Errorf("failed to create ports: %w", err)
	}
	binds, err := c.createBinds()
	if err != nil {
		return fmt.Errorf("failed to create binds: %w", err)
	}
	image := c.Image
	if _, err := os.Stat(filepath.Join(c.Image, "Dockerfile")); err == nil {
		image = c.name
	}

	log.Printf("creating container")
	_, err = cli.ContainerCreate(ctx, &dockercontainer.Config{
		Hostname:     c.name,
		ExposedPorts: portSet,
		Tty:          c.TTY,
		Env:          environ,
		Cmd:          strslice.StrSlice(c.Args),
		Image:        image,
		User:         c.User,
		WorkingDir:   c.WorkingDir,
		Entrypoint:   strslice.StrSlice(c.GetCommand()),
		Labels:       map[string]string{hashLabel: expectedHash},
	}, &dockercontainer.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
	}, &network.NetworkingConfig{}, &v1.Platform{}, c.name)
	if ignoreConflict(err) != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	id, _, err = c.getContainer(ctx, cli)
	if err != nil {
		return fmt.Errorf("failed to get container ID: %w", err)
	}

	c.containerID = id
	if err = cli.ContainerStart(ctx, id, dockertypes.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	go func() {
		<-ctx.Done()
		if err := c.stop(context.Background()); err != nil {
			log.Printf("failed to stop: %v", err)
		}
	}()
	logs, err := cli.ContainerLogs(ctx, c.name, dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Since:      time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to log container: %w", err)
	}
	defer logs.Close()
	if _, err = stdcopy.StdCopy(stdout, stderr, logs); err != nil {
		// ignore errors, might be content cancelled, we still need to wait for the container to exit
		log.Printf("failed to log container: %v", err)
	}

	waitC, errC := cli.ContainerWait(context.Background(), id, dockercontainer.WaitConditionNotRunning)
	select {
	case wait := <-waitC:
		if wait.StatusCode != 0 {
			return fmt.Errorf("exit code %d", wait.StatusCode)
		}
		return nil
	case err := <-errC:
		return fmt.Errorf("failed to wait for container: %w", err)
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
		for _, volume := range c.spec.Volumes {
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

func (c *container) stop(ctx context.Context) error {
	if c.name == "" {
		return nil
	}
	log := c.log
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()
	id, _, err := c.getContainer(ctx, cli)
	if err != nil {
		return fmt.Errorf("failed to get container ID: %w", err)
	}
	if id == "" {
		return nil
	}
	log.Printf("stopping container\n")
	grace := c.spec.GetTerminationGracePeriod()
	timeout := int(grace.Seconds())
	err = cli.ContainerStop(ctx, id, dockercontainer.StopOptions{
		Timeout: &timeout,
	})
	if ignoreNotExist(err) != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

const hashLabel = "kit.hash"

func (c *container) getContainer(ctx context.Context, cli *client.Client) (string, string, error) {
	list, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{All: true})
	if err != nil {
		return "", "", err
	}
	for _, existing := range list {
		if slices.Contains(existing.Names, "/"+c.name) {
			id := existing.ID
			return id, existing.Labels[hashLabel], nil
		}
	}
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

func (c *container) GetMetrics(ctx context.Context) (*types.Metrics, error) {
	if c.containerID == "" {
		return &types.Metrics{}, nil
	}

	command := metrics.GetProcFSCommand(1) // PID 1
	cmdArgs := append([]string{"exec", c.name}, command...)
	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker exec ps failed for container %s: %w", c.name, err)
	}

	metrics, procFSSnapshot, err := metrics.ParseProcFSOutput(string(output), c.profFSSnapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to parse process metrics for container %s: %w", c.name, err)
	}
	c.profFSSnapshot = procFSSnapshot
	return metrics, nil
}

var _ Interface = &container{}
