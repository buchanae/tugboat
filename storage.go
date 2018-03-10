package tugboat

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

type Storage interface {
	Get(ctx context.Context, url, rel, abs string) error
	Put(ctx context.Context, url, rel, abs string) error

	// Determines whether this backends supports the given request (url/path/class).
	// A backend normally uses this to match the url prefix (e.g. "s3://")
	SupportsGet(url string) bool
	SupportsPut(url string) bool
}

func Download(ctx context.Context, store Storage, log Logger, inputs []File) error {
	return nil
}

func Upload(ctx context.Context, stage *StagedTask, store Storage, log Logger) error {

	errors := make(chan error)
	files := make(chan *hostfile)
	done := make(chan struct{})
	wg := &sync.WaitGroup{}

	// Start a fixed number of uploader threads.
	numUploaders := 10
	wg.Add(numUploaders)
	for i := 0; i < numUploaders; i++ {
		go func() {
			defer wg.Done()

			for file := range files {
				log.UploadStarted(file.out)

				// TODO
				//r.fixLinks(mapper, output.Path)
				// TODO log bytes copied
				rel := stage.Unmap(file.path)
				err := store.Put(ctx, file.out.URL, rel, file.path)
				if err != nil {
					errors <- wrap(err, "upload failed %s, %s", file.out.URL, file.path)
				} else {
					log.UploadFinished(file.out)
				}
			}
		}()
	}

	// Collect errors
	var me MultiError
	go func() {
		for err := range errors {
			me = append(me, err)
		}
		close(done)
	}()

	// Walk all the outputs, sending files to the uploader channel.
	for _, out := range stage.Outputs {
		w := walker{out, files, errors}
		filepath.Walk(out.Path, w.walk)
	}

	close(files)
	wg.Wait()
	close(errors)
	<-done

	return me.Finish()
}

type hostfile struct {
	out File
	// The absolute path of the file on the host.
	path string
	// Size of the file in bytes
	size int64
}

type walker struct {
	out   File
	files chan *hostfile
	errs  chan error
}

func (w *walker) walk(p string, f os.FileInfo, err error) error {
	if err != nil {
		w.errs <- err
		// Skip this file/directory, capture the error, and continue processing.
		return nil
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		w.errs <- err
		// Skip this file/directory, capture the error, and continue processing.
		return nil
	}

	if !f.IsDir() {
		w.files <- &hostfile{w.out, abs, f.Size()}
	}
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
