package gs

import (
	"context"
	"fmt"
	"io"
	"os"

  tug "github.com/buchanae/tugboat"
  "cloud.google.com/go/storage"
)

type GS struct {
	workdir string
	svc *storage.Client
}

func NewGS(workdir string) (*GS, error) {
	ctx := context.Background()
  client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GS{workdir, client}, nil
}

func (gs *GS) Get(ctx context.Context, rawurl string, hostPath string) error {

  // TODO directory
  err := tug.EnsurePath(hostPath, 0755)
  if err != nil {
    return err
  }

  bkt := gs.svc.Bucket(gs.workdir)
  obj := bkt.Object(rawurl)
  reader, err := obj.NewReader(ctx)
  if err != nil {
    return err
  }

  fh, err := os.Create(hostPath)
  if err != nil {
    return err
  }
  defer fh.Close()

  _, err = io.Copy(fh, reader)
  if err != nil {
    return err
  }
  return nil
}

// PutFile copies an object (file) from the host path to GS.
func (gs *GS) Put(ctx context.Context, rawurl, rel, hostPath string) error {
  return fmt.Errorf("put unimplemented")
}

func (gs *GS) SupportsGet(rawurl string) bool {
  return true
}

func (gs *GS) SupportsPut(rawurl string) bool {
  return true
}
