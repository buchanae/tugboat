package tugboat

import (
	"os"
	"path/filepath"
	"strings"
)

func StageTask(stage *Stage, task *Task) (*Task, error) {
	return nil, nil
}

type Stage struct {
	Dir  string
	Mode os.FileMode
}

func NewStage(dir string, mode os.FileMode) (*Stage, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, wrap(err, "failed to get absolute path")
	}

	err = EnsureDir(dir, mode)
	if err != nil {
		return nil, wrap(err, "failed to create stage directory")
	}

	return &Stage{dir, mode}, nil
}

// Map maps the given path into the stage directory.
// An error is returned if the resulting path would be outside the stage directory.
//
// If the stage is configured with a base dir of "/tmp/staging", then
// stage.Map("/home/ubuntu/myfile") will return "/tmp/staging/home/ubuntu/myfile".
func (stage *Stage) Map(src string) (string, error) {
	p := filepath.Join(stage.Dir, src)
	ap, err := filepath.Abs(p)
	if err != nil {
		return "", wrap(err, "failed to get absolute path")
	}
	if !strings.HasPrefix(ap, stage.Dir) {
		return "", errf("invalid path: %s is not a valid subpath of %s", p, stage.Dir)
	}
	return ap, nil
}

// Unmap strips the stage directory prefix from the given path.
//
// If the stage is configured with a base dir of "/tmp/staging", then
// stage.Unmap("/tmp/staging/home/ubuntu/myfile") will return "/home/ubuntu/myfile".
func (stage *Stage) Unmap(src string) string {
	p := strings.TrimPrefix(src, stage.Dir)
	p = filepath.Clean("/" + p)
	return p
}

// RemoveAll removes the stage directory.
func (stage *Stage) RemoveAll() error {
	return os.RemoveAll(stage.Dir)
}

// exists returns whether the given file or directory exists or not
func exists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, wrap(err, "failed to call os.Stat")
}

// EnsureDir ensures a directory exists.
func EnsureDir(p string, mode os.FileMode) error {
	s, err := os.Stat(p)
	if err == nil {
		if s.IsDir() {
			return nil
		}
		return errf("file exists but is not a directory")
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(p, mode)
		if err != nil {
			return wrap(err, "failed to create directory")
		}
		return nil
	}
	return err
}

// EnsurePath ensures a directory exists, given a file path. This calls path.Dir(p)
func EnsurePath(p string, mode os.FileMode) error {
	dir := filepath.Dir(p)
	return EnsureDir(dir, mode)
}
