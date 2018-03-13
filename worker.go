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

	Workdir string

	Volumes []string
	Inputs,
	// All output paths must be contained in a volume.
	Outputs []File

	Stdin, Stdout, Stderr string
}

type Executor interface {
	Exec(context.Context, *StagedTask, *Stdio) error
}

func Run(ctx context.Context, task *Task, stage *Stage, log Logger, store Storage, exec Executor) (err error) {

	var me MultiError
	try := me.Try
	defer func() {
		err = me.Finish()
	}()

	info := log.Info
	d := LogHelper{log}
	d.Start()
	defer d.Finish()

	// TODO
	//info("validating task")
	//err = store.Validate(ctx, task.Outputs)
	//Must(err)

	info("creating staging directory")
	var staged *StagedTask
	staged, err = StageTask(stage, task)
	try(err)
	if err != nil {
		return
	}

	defer func() {
		try(staged.RemoveAll())
	}()

	err = Download(ctx, staged, store, log)
	try(err)
	if err != nil {
		return
	}

	defer func() {
		try(Upload(ctx, staged, store, log))
	}()

	var stdio *Stdio
	stdio, err = DefaultStdio(staged, log)
	try(err)
	if err != nil {
		return
	}

	defer func() {
		try(stdio.Close())
	}()
	defer info("cleaning up")

	log.Running()
	try(exec.Exec(ctx, staged, stdio))

	return
}
