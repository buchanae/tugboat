package tugboat

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

type Storage interface {
	Get(ctx context.Context, url, abs string) error
	Put(ctx context.Context, url, rel, abs string) error

	// Determines whether this backends supports the given request (url/path/class).
	// A backend normally uses this to match the url prefix (e.g. "s3://")
	SupportsGet(url string) bool
	SupportsPut(url string) bool
}

func Download(ctx context.Context, task *StagedTask, store Storage, log Logger) error {

	errors := make(chan error)
	files := make(chan File)
	done := make(chan struct{})
	wg := &sync.WaitGroup{}

	// Start a fixed number of downloader threads.
	numDownloaders := 10
	wg.Add(numDownloaders)
	for i := 0; i < numDownloaders; i++ {
		go func() {
			defer wg.Done()

			for file := range files {
				log.DownloadStarted(file)
				// TODO log bytes copied

				err := store.Get(ctx, file.URL, file.Path)
				if err != nil {
					errors <- wrap(err, "download failed %s, %s", file.URL, file.Path)
				} else {
					log.DownloadFinished(file)
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

	for _, input := range task.Inputs {
		files <- input
	}
	close(files)

	wg.Wait()
	close(errors)
	<-done

	return me.Finish()
}

func Upload(ctx context.Context, task *StagedTask, store Storage, log Logger) error {

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

				err := store.Put(ctx, file.out.URL, file.rel, file.path)
				if err != nil {
					errors <- wrap(err, "uploading %q to %q", file.path, file.out.URL)
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
	for _, out := range task.Outputs {
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
	rel string
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

	rel, err := filepath.Rel(w.out.Path, p)
	if err != nil {
		w.errs <- wrap(err, "getting relative path for %q", p)
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
		w.files <- &hostfile{w.out, rel, abs, f.Size()}
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

*/
