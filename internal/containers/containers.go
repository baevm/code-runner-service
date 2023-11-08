package containers

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	JS_IMAGE_KEY = "javascript"
	PY_IMAGE_KEY = "python"
)

var images = map[string]string{
	JS_IMAGE_KEY: "node:20-alpine",
	PY_IMAGE_KEY: "python:3.11-slim-bookworm",
}

func PullImages() error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())

	if err != nil {
		return err
	}

	for _, refStr := range images {
		_, err := cli.ImagePull(ctx, refStr, types.ImagePullOptions{})

		if err != nil {
			return err
		}
	}

	return nil
}

func RunCodeContainer(lang string, code string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	defer reader.Close()

	if _, isExist := images[lang]; !isExist {
		return "", errors.New("language not found")
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: images[lang],
		Cmd:   []string{"python3 ", code},
		Tty:   true,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}

	case <-statusCh:

	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)

	_, err = io.Copy(buf, out)
	if err != nil {
		return "", err
	}

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
