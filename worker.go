package tugboat

import (
	"context"
)

type SystemError struct{}
type ExecError struct {
	ExitCode int
}

type InvalidInputsError struct{}
type InvalidOutputsError struct{}

type File struct {
	URL  string
	Path string
}

type Task struct {
	ID             string
	ContainerImage string
	Command        []string
	Env            map[string]string

	Volumes []string
	Inputs,
	// All output paths must be contained in a volume.
	Outputs []File

	Stdin, Stdout, Stderr string
}

type Executor interface {
	Exec(context.Context, *Task, *Stdio) error
}

type EmptyExecutor struct{}

func (e *EmptyExecutor) Exec(ctx context.Context, task *Task, stdio *Stdio) error {
	return nil
}

func Run(ctx context.Context, task *Task, log Logger, store Storage, exec Executor) (err error) {

	try, must, finish := Errors()
	defer func() { err = finish(err) }()

	info := log.Info
	d := LogHelper{log}
	d.Start()
	defer d.Finish()

	info("validating task")
	must(store.Validate(ctx, task.Outputs))

	info("creating staging directory")
	stage, err := NewStage("foo", 0755)
	must(err)

	staged, err := StageTask(stage, task)
	must(err)

	defer func() {
		info("cleaning up staging directory")
		try(stage.RemoveAll())
	}()

	info("downloading inputs")
	must(store.Download(ctx, staged.Inputs))

	defer func() {
		info("uploading outputs")
		try(store.Upload(ctx, staged.Outputs))
	}()

	log.Running()

	info("opening stdio")
	stdio, err := DefaultStdio(staged, log)
	must(err)

	defer func() {
		info("closing stdio")
		try(stdio.Close())
	}()

	must(exec.Exec(ctx, task, stdio))
	return
}
