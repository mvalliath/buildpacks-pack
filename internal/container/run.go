package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dcontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"io"
	"time"
)

func Run(ctx context.Context, docker client.CommonAPIClient, ctrID string, out, errOut io.Writer) error {
	bodyChan, errChan := docker.ContainerWait(ctx, ctrID, dcontainer.WaitConditionNextExit)

	if err := docker.ContainerStart(ctx, ctrID, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "container start")
	}

	//wait for container to start or logs won't
	time.Sleep(10*time.Millisecond)

	logs, err := docker.ContainerLogs(ctx, ctrID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return errors.Wrap(err, "container logs stdout")
	}
	if logs != nil {
		defer logs.Close()
	}

	copyErr := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(out, errOut, logs)
		copyErr <- err
	}()

	select {
	case body := <-bodyChan:
		if body.StatusCode != 0 {
			return fmt.Errorf("failed with status code: %d", body.StatusCode)
		}
	case err := <-errChan:
		copyErr <- err
	}

	return <-copyErr
}
