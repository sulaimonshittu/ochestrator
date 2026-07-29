package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
	Runtime       Task
}

type Docker struct {
	Client *client.Client
	Config
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(
		ctx,
		d.Image,
		image.PullOptions{},
	)
	if err != nil {
		log.Printf("Error pulling image %s: error: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.RestartPolicy),
	}
	r := container.Resources{
		Memory:   d.Memory,
		NanoCPUs: int64(d.Cpu * math.Pow(10, 9)),
	}
	cc := container.Config{
		Image:        d.Image,
		Tty:          false,
		Env:          d.Env,
		ExposedPorts: d.ExposedPorts,
	}
	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Name)
	if err != nil {
		log.Printf("Error creating ccontainer using image %s: error: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: error %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	d.Runtime.ContainerID = resp.ID
	out, err := d.Client.ContainerLogs(
		ctx,
		resp.ID,
		container.LogsOptions{},
	)
	if err != nil {
		log.Printf("Error getting logs for container %s: error: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return DockerResult{ContainerId: resp.ID, Action: "Start", Result: "Success"}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: error: %v\n", id, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerRemove(ctx, id, container.RemoveOptions{
		true,
		false,
		false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}
	return DockerResult{
		Action: "stop",
		Result: "success",
		Error:  nil,
	}
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}
type Task struct {
	ID          uuid.UUID
	ContainerID string
	Name        string
	State
	Image         string
	CPU           float64
	Memory        int64
	Disk          int64
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID uuid.UUID
	State
	Timestamp time.Time
	Task
}
