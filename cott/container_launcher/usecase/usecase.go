package usecase

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

var STOP_CONTAINER_TIMEOUT = 10 * time.Second

type ContainerLauncherUsecase interface {
	// Start continer and returns container ID on success
	LaunchContainer(image string, envVarMap map[string]string, port uint16) (*string, error)
	StopContainer(id string) error
	RemoveContainer(id string) error
	// GetContainerStats get channel with container stats and cancel func for stopping receiving container stats
	GetContainerStats(id string) (*types.StatsJSON, error)
	GetContainerStatsStream(id string) (<-chan *types.Stats, context.CancelFunc, error)
}

type containerLauncherUsecase struct {
	cli *client.Client
}

func NewContainerLauncherUsecase() (ContainerLauncherUsecase, error) {
	cluc := new(containerLauncherUsecase)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cluc.cli = cli

	return cluc, nil
}

func (cluc *containerLauncherUsecase) LaunchContainer(image string, envVarMap map[string]string, port uint16) (*string, error) {
	logrus.WithFields(logrus.Fields{"image": image, "envVarMap": envVarMap, "port": port}).Debug("launch container")

	if reader, err := cluc.cli.ImagePull(context.Background(), image, types.ImagePullOptions{}); err != nil {
		return nil, err
	} else {
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, reader); err != nil {
			return nil, err
		}
		logrus.Trace(buf)
		reader.Close()
	}
	logrus.WithFields(logrus.Fields{"image": image}).Debug("container image pulled")

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

	resp, err := cluc.cli.ContainerCreate(context.Background(), containerCfg, hostCfg, nil, nil, "")
	if err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"image": image, "id": resp.ID}).Debug("container created")

	if err := cluc.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	logrus.WithFields(logrus.Fields{"image": image, "id": resp.ID}).Debug("container started")

	return &resp.ID, nil
}

func (cluc *containerLauncherUsecase) StopContainer(id string) error {
	if err := cluc.cli.ContainerStop(context.Background(), id, &STOP_CONTAINER_TIMEOUT); err != nil {
		return err
	}
	logrus.WithField("id", id).Debug("container stopped")

	return nil
}

func (cluc *containerLauncherUsecase) RemoveContainer(id string) error {
	if err := cluc.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	logrus.WithField("id", id).Debug("container removed")

	return nil
}

func (cluc *containerLauncherUsecase) GetContainerStats(id string) (*types.StatsJSON, error) {
	statsResponse, err := cluc.cli.ContainerStats(context.Background(), id, false)
	if err != nil {
		return nil, err
	}

	statsBytes, err := io.ReadAll(statsResponse.Body)
	if err != nil {
		return nil, err
	}

	var stats types.StatsJSON

	if err := json.Unmarshal(statsBytes, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (cluc *containerLauncherUsecase) GetContainerStatsStream(id string) (<-chan *types.Stats, context.CancelFunc, error) {
	ctx, ctxCancelFunc := context.WithCancel(context.Background())

	statsResponse, err := cluc.cli.ContainerStats(ctx, id, true)
	if err != nil {
		ctxCancelFunc()
		return nil, nil, err
	}
	logrus.WithField("id", id).Debug("start getting container stats")

	statsCh := make(chan *types.Stats, 1)

	// Goroutine for sending data from stats to channel
	go func() {
		defer close(statsCh)

		decoder := json.NewDecoder(statsResponse.Body)
		for {
			select {
			case <-ctx.Done():
				statsResponse.Body.Close()
				logrus.WithField("id", id).Debug("stop logging container stats. context done")
				return
			default:
				var stats types.Stats
				if err := decoder.Decode(&stats); err == io.EOF {
					logrus.WithField("id", id).Debug("stop logging container stats")
					return
				} else if err != nil {
					ctxCancelFunc()
					break
				}
				statsCh <- &stats
			}
		}
	}()

	return statsCh, ctxCancelFunc, nil
}

func (cluc *containerLauncherUsecase) convertEnvVarsMapToSlice(envVarMap map[string]string) []string {
	var envVarsSlice []string
	for k, v := range envVarMap {
		envVarsSlice = append(envVarsSlice, k+"="+v)
	}
	return envVarsSlice
}
