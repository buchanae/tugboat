package tugboat

import (
	"io"
	"io/ioutil"
	"os"
)

type Stdio struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

func (s *Stdio) Close() error {
	var errors MultiError
	if closer, ok := s.Stdout.(io.WriteCloser); ok {
		err := closer.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	if closer, ok := s.Stderr.(io.WriteCloser); ok {
		err := closer.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func DefaultStdio(task *Task, log Logger) (*Stdio, error) {
	stdio, err := TaskStdio(task)
	if err != nil {
		return nil, err
	}
	LogStdio(log, stdio)
	return stdio, nil
}

func TaskStdio(task *Task) (*Stdio, error) {
	stdio := &Stdio{
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}

	if task.Stdin != "" {
		s, err := os.Open(task.Stdin)
		if err != nil {
			return nil, wrap(err, "failed to open stdin file")
		}
		stdio.Stdin = s
	}

	if task.Stdout != "" {
		s, err := os.Create(task.Stdout)
		if err != nil {
			return nil, wrap(err, "failed to create stdout file")
		}
		stdio.Stdout = s
	}

	if task.Stderr != "" {
		s, err := os.Create(task.Stderr)
		if err != nil {
			return nil, wrap(err, "failed to create stderr file")
		}
		stdio.Stderr = s
	}
	return stdio, nil
}

func LogStdio(log Logger, stdio *Stdio) {
	stdout := log.Stdout()
	stderr := log.Stderr()

	if stdout != nil {
		stdio.Stdout = io.MultiWriter(stdio.Stdout, stdout)
	}
	if stderr != nil {
		stdio.Stderr = io.MultiWriter(stdio.Stderr, stderr)
	}
}
