package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

var (
	STOP_CONTAINER_TIMEOUT = 10 * time.Second
)

type ContainerLauncherUsecase interface {
	// Start continer and returns container ID on success
	LaunchContainer(image string, name string, envVarMap map[string]string, port uint16) (*string, error)
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

func (cluc *containerLauncherUsecase) LaunchContainer(image string, name string, envVarMap map[string]string, port uint16) (*string, error) {
	logrus.WithFields(logrus.Fields{"image": image, "name": name, "envVarMap": envVarMap, "port": port}).Debug("launch container")

	reader, err := cluc.cli.ImagePull(cluc.ctx, image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"image": image}).Debug("container image pulled")

	defer reader.Close()

	portStr := strconv.FormatUint(uint64(port), 10)
	containerPort := nat.Port(portStr)
	containerCfg := &container.Config{
		Image: image,
		Env:   cluc.convertEnvVarsMapToSlice(envVarMap),
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
		},
	}
	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				nat.PortBinding{
					HostIP:   "0.0.0.0",
					HostPort: portStr,
				},
			},
		},
	}

	resp, err := cluc.cli.ContainerCreate(cluc.ctx, containerCfg, hostCfg, nil, nil, name)
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"image": image, "name": name, "id": resp.ID}).Debug("container created")

	if err := cluc.cli.ContainerStart(cluc.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"image": image, "name": name, "id": resp.ID}).Debug("container started")

	return &resp.ID, nil
}

func (cluc *containerLauncherUsecase) StopContainer(id string) error {
	if err := cluc.cli.ContainerStop(cluc.ctx, id, &STOP_CONTAINER_TIMEOUT); err != nil {
		return err
	}
	logrus.WithField("id", id).Debug("container stopped")

	return nil
}

func (cluc *containerLauncherUsecase) RemoveContainer(id string) error {
	if err := cluc.cli.ContainerRemove(cluc.ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	logrus.WithField("id", id).Debug("container removed")

	return nil
}

func (cluc *containerLauncherUsecase) convertEnvVarsMapToSlice(envVarMap map[string]string) []string {
	var envVarsSlice []string
	for k, v := range envVarMap {
		envVarsSlice = append(envVarsSlice, k+"="+v)
	}
	return envVarsSlice
}
