package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local provides access to a local-disk storage system.
type Local struct {}

// NewLocal returns a Local instance.
func NewLocal() (*Local, error) {
	return &Local{}, nil
}

// Get copies a file from storage into the given hostPath.
func (local *Local) Get(ctx context.Context, url, host string) error {
	path := getPath(url)
  return linkFile(path, host)
}

// Put copies a file from the path into the storage url.
func (local *Local) Put(ctx context.Context, url, rel, host string) error {
	path := getPath(url)
  tgt := filepath.Join(path, rel)
	return linkFile(host, tgt)
}

// SupportsGet returns true if the Local storage driver is able to get this file.
func (local *Local) SupportsGet(url string) bool {
	if !strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "file://") {
    return false
	}
  return true
}

// SupportsPut returns true if the Local storage driver is able to put this file.
func (local *Local) SupportsPut(url string) bool {
	return local.SupportsGet(url)
}

func getPath(rawurl string) string {
	p := strings.TrimPrefix(rawurl, "file://")
	return p
}

// Copies file source to destination dest.
func copyFile(source string, dest string) (err error) {
	// check if dest exists; if it does check if it is the same as the source
	same, err := sameFile(source, dest)
	if err != nil {
		return err
	}
	if same {
		return nil
	}
	// Open source file for copying
	sf, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file for copying: %s", err)
	}
	defer sf.Close()

	// Create and open dest file for writing
	df, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return fmt.Errorf("failed to create dest file for copying: %s", err)
	}
	defer func() {
		cerr := df.Close()
		if cerr != nil {
			err = fmt.Errorf("%v; %v", err, cerr)
		}
	}()

	_, err = io.Copy(df, sf)
	return err
}

// Hard links file source to destination dest.
func linkFile(source string, dest string) error {
	var err error
	// without this resulting link could be a symlink
	parent, err := filepath.EvalSymlinks(source)
	if err != nil {
		return fmt.Errorf("failed to eval symlinks: %s", err)
	}
	same, err := sameFile(parent, dest)
	if err != nil {
		return fmt.Errorf("failed to check if file is the same file: %s", err)
	}
	if same {
		return nil
	}
	err = os.Link(parent, dest)
	if err != nil {
		err = copyFile(source, dest)
		if err != nil {
			return fmt.Errorf("failed to copy file: %s", err)
		}
	}
	return err
}

func sameFile(source string, dest string) (bool, error) {
	var err error
	sfi, err := os.Stat(source)
	if err != nil {
		return false, fmt.Errorf("failed to stat src file: %s", err)
	}
	dfi, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to stat dest file: %s", err)
	}
	return os.SameFile(sfi, dfi), nil
}
