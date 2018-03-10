package docker

import (
	"context"
	"encoding/json"
	"fmt"
	tug "github.com/buchanae/tugboat"
	"os/exec"
	"strings"
	"time"
)

type Docker struct {
	tug.Logger
	LeaveContainer bool
}

func (d *Docker) Exec(ctx context.Context, task *tug.StagedTask, stdio *tug.Stdio) error {

	pullErr := exec.Command("docker", "pull", task.ContainerImage).Run()
	if pullErr != nil {
		d.Info(`failed to pull container image "%s": %s`, task.ContainerImage, pullErr)
	}

	args := []string{"run", "-i", "--read-only"}

	if !d.LeaveContainer {
		args = append(args, "--rm")
	}

	for k, v := range task.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	name := fmt.Sprintf("task-%s-%s", task.ID, randString(5))
	args = append(args, "--name", name)

	for i, input := range task.Inputs {
		host := input.Path
		container := task.Task.Inputs[i].Path
		arg := formatVolumeArg(host, container, true)
		args = append(args, "-v", arg)
	}

	for i, host := range task.Volumes {
		container := task.Task.Volumes[i]
		arg := formatVolumeArg(host, container, false)
		args = append(args, "-v", arg)
	}

	args = append(args, task.ContainerImage)
	args = append(args, task.Command...)

	// Roughly: `docker run --rm -i --read-only -w [workdir] -v [bindings] [imageName] [cmd]`
	d.Meta("command", "docker "+strings.Join(args, " "))
	d.Meta("container name", name)

	cmd := exec.Command("docker", args...)

	cmd.Stdin = stdio.Stdin
	cmd.Stdout = stdio.Stdout
	cmd.Stderr = stdio.Stderr

	var err error

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf(`exec "docker run" failed: %s`, err)
	}

	cmdctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Kill the container when the context is canceled,
	// instead of expecting the os/exec signal to work.
	go func() {
		<-cmdctx.Done()
		exec.Command("docker", "stop", "-t", "10", name).Run()
		exec.Command("docker", "kill", name).Run()
	}()

	// Inspect the container for metadata
	go func() {
		ticker := time.NewTicker(time.Second)
		cmd := exec.CommandContext(cmdctx, "docker", "inspect", name)
		for i := 0; i < 5; i++ {
			select {
			case <-cmdctx.Done():
				return
			case <-ticker.C:
				out, err := cmd.Output()
				if err == nil {
					meta := ContainerMetadata{}
					err := json.Unmarshal(out, &meta)
					if err == nil {
						d.Meta("container ID", meta.Id)
						d.Meta("container image hash", meta.Image)
						return
					}
				}
			}
		}
	}()

	return cmd.Wait()
}

type ContainerMetadata struct {
	Id    string
	Image string
}
