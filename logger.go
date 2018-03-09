package tugboat

import (
	"io"
	"os"
	"time"
)

type Logger interface {
	StartTime(time.Time)
	EndTime(time.Time)

	Meta(key string, value interface{})
	Version(Version)

	Info(msg string, args ...interface{})

	DownloadStarted(url string)
	DownloadFinished(url string)
	UploadStarted(url string)
	UploadFinished(url string)

	Running()

	Stdout() io.Writer
	Stderr() io.Writer
}

type EmptyLogger struct {
	Logger
}

type LogHelper struct {
	Logger
}

func (d *LogHelper) Start() {
	d.Logger.Version(version)

	if name, err := os.Hostname(); err == nil {
		d.Logger.Meta("hostname", name)
	}
	d.Logger.StartTime(time.Now())
}

func (d *LogHelper) Finish() {
	d.Logger.EndTime(time.Now())
}

type Version struct {
	Name, Doc                string
	Major, Minor, Patch      int
	Commit, Branch, Upstream string
	Date                     time.Time
}

var version = Version{}
