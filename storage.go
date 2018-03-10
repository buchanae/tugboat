package tugboat

import (
	"context"
	"fmt"
)

type Storage interface {
	Validate(context.Context, []File) error
	Download(context.Context, []File) error
	Upload(context.Context, []File) error
}

type EmptyStorage struct {
}

func (s *EmptyStorage) Validate(ctx context.Context, files []File) error {
	fmt.Println("Validate", files)
	return nil
}
func (s *EmptyStorage) Download(ctx context.Context, files []File) error {
	fmt.Println("Download", files)
	return nil
}
func (s *EmptyStorage) Upload(ctx context.Context, files []File) error {
	fmt.Println("Upload", files)
	return nil
}

/*
// Validate the output uploads
func validateOutputs(mapper *FileMapper) error {
	for _, output := range mapper.Outputs {
		err := r.Store.SupportsPut(output.Url, output.Type)
		if err != nil {
			return fmt.Errorf("Output upload not supported by storage: %v", err)
		}
	}
	return nil
}


	// Download inputs
	downloadErrs := make(util.MultiError, len(mapper.Inputs))
	downloadCtx, cancelDownloadCtx := context.WithCancel(ctx)
	defer cancelDownloadCtx()

	wg := &sync.WaitGroup{}

  for i, input := range mapper.Inputs {
    wg.Add(1)

    go func(input *tes.Input, i int) {
      defer wg.Done()

      log.DownloadStarted(input.URL)

      err := r.Store.Get(downloadCtx, input.Url, input.Path, input.Type)
      if err != nil {
        downloadErrs[i] = err
        event.Error("Download failed", "url", input.Url, "error", err)
        cancelDownloadCtx()
      } else {
        log.DownloadFinished(input.URL)
      }

    }(input, i)
  }

	wg.Wait()
	if !downloadErrs.IsNil() {
		run.syserr = downloadErrs.ToError()
	}

	// Upload outputs
	uploadErrs := make(util.MultiError, len(mapper.Outputs))
	wg = &sync.WaitGroup{}

  for i, output := range mapper.Outputs {
    wg.Add(1)
    go func(output *tes.Output, i int) {
      defer wg.Done()

      // TODO map files to URLs outside of storage.Put
      log.UploadStarted(output.URL)

      r.fixLinks(mapper, output.Path)

      out, err := r.Store.Put(ctx, output.Url, output.Path, output.Type)
      if err != nil {
        if err == storage.ErrEmptyDirectory {
          event.Warn("Upload finished with warning", "url", output.Url, "warning", err)
        } else {
          uploadErrs[i] = err
          event.Error("Upload failed", "url", output.Url, "error", err)
        }
      } else {
        log.UploadFinished(output.URL)
      }
      outputs[i] = out
      return
    }(output, i)
  }

	wg.Wait()

	if !uploadErrs.IsNil() {
		run.syserr = uploadErrs.ToError()
	}

	// unmap paths for OutputFileLog
	for _, o := range outputLog {
		o.Path = mapper.ContainerPath(o.Path)
	}
*/
