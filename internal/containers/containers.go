package containers

import (
	"code-runner-service/internal/models"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	JS_IMAGE_KEY = "javascript"
	PY_IMAGE_KEY = "python"
)

var images = map[string]string{
	JS_IMAGE_KEY: "node:20-alpine",
	PY_IMAGE_KEY: "python:3.11-slim-bookworm",
}

var languagePrefix = map[string]string{
	JS_IMAGE_KEY: "js",
	PY_IMAGE_KEY: "py",
}

type Client struct {
	cli *client.Client
}

func New() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
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

// Creates temp file, writes to it.
func createFile(code models.Code) (*os.File, error) {
	file, err := os.CreateTemp("./temp", fmt.Sprintf("%s.*.%s", code.Lang, languagePrefix[code.Lang]))

	if err != nil {
		return nil, err
	}

	codeData := []byte(code.Body)

	if _, err := file.Write(codeData); err != nil {
		return nil, err
	}

	return file, err
}

func (c *Client) RunCodeContainer(ctx context.Context, code models.Code) (string, error) {
	if _, isExist := images[code.Lang]; !isExist {
		return "", errors.New("language not found")
	}

	//ctx := context.Background()

	file, err := createFile(code)

	if err != nil {
		return "", err
	}

	defer file.Close()
	defer os.Remove(file.Name())

	codeArchive, err := archive.Tar(file.Name(), archive.Gzip)

	if err != nil {
		return "", err
	}

	defer codeArchive.Close()

	fileFields := strings.Split(file.Name(), "/")
	fileName := fileFields[len(fileFields)-1]

	cont, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image:        images[code.Lang],
		Tty:          true,
		AttachStdout: true,
		Cmd:          []string{"bash", "-c", fmt.Sprintf("python3 %s", fileName)},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	defer c.cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})

	err = c.cli.CopyToContainer(ctx, cont.ID, "/", codeArchive, types.CopyToContainerOptions{})

	if err != nil {
		return "", err
	}

	if err := c.cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := c.cli.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}

	case <-statusCh:

	}
	out, err := c.cli.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
	})
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)

	// https://stackoverflow.com/questions/52774830/docker-exec-command-from-golang-api
	_, err = io.Copy(buf, out)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
