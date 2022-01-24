package usecase

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	STOP_CONTAINER_TIMEOUT = 10 * time.Second
)

type ContainerLauncherUsecase interface {
	// Start continer and returns container ID on success
	LaunchContainer(image string, name string, ports []string) (*string, error)
	StopContainer(id string) error
	RemoveContainer(id string) error
}

type containerLauncherUsecase struct {
	ctx context.Context
	cli *client.Client
}

func NewContainerLauncherUsecase() (ContainerLauncherUsecase, error) {
	cluc := new(containerLauncherUsecase)
	cluc.ctx = context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	cluc.cli = cli

	return cluc, nil
}

func (cluc *containerLauncherUsecase) LaunchContainer(image string, name string, ports []string) (*string, error) {
	reader, err := cluc.cli.ImagePull(cluc.ctx, image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	containerCfg := &container.Config{
		Hostname: name,
		Image:    image}

	for _, port := range ports {
		containerCfg.ExposedPorts[nat.Port(port)] = struct{}{}
	}

	resp, err := cluc.cli.ContainerCreate(cluc.ctx, containerCfg, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cluc.cli.ContainerStart(cluc.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := cluc.cli.ContainerWait(cluc.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	return &resp.ID, nil
}

func (cluc *containerLauncherUsecase) StopContainer(id string) error {
	if err := cluc.cli.ContainerStop(cluc.ctx, id, &STOP_CONTAINER_TIMEOUT); err != nil {
		return err
	}

	return nil
}

func (cluc *containerLauncherUsecase) RemoveContainer(id string) error {
	if err := cluc.cli.ContainerRemove(cluc.ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}
