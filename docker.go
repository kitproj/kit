package main

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"path/filepath"
)

func build(ctx context.Context, name string, cli *client.Client, path string, stdout io.Writer) error {

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerfileName := filepath.Base(path)
	readDockerFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tarHeader := &tar.Header{
		Name: dockerfileName,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return err
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerfileName,
			Remove:     true,
			Tags:       []string{name},
		})
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}
	return err
}
