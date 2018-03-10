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

func DefaultStdio(t *StagedTask, log Logger) (*Stdio, error) {
	stdio, err := NewStdio(t.Stdin, t.Stdout, t.Stderr)
	if err != nil {
		return nil, err
	}
	LogStdio(log, stdio)
	return stdio, nil
}

func NewStdio(stdin, stdout, stderr string) (*Stdio, error) {
	stdio := &Stdio{
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}

	if stdin != "" {
		s, err := os.Open(stdin)
		if err != nil {
			return nil, wrap(err, "failed to open stdin file")
		}
		stdio.Stdin = s
	}

	if stdout != "" {
		s, err := os.Create(stdout)
		if err != nil {
			return nil, wrap(err, "failed to create stdout file")
		}
		stdio.Stdout = s
	}

	if stderr != "" {
		s, err := os.Create(stderr)
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
