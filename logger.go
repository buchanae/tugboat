package tugboat

import (
	"fmt"
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
}

func (e EmptyLogger) StartTime(t time.Time) {
	fmt.Println("StartTime", t)
}
func (e EmptyLogger) EndTime(t time.Time) {
	fmt.Println("EndTime", t)
}
func (e EmptyLogger) Meta(key string, value interface{}) {
	fmt.Println("Meta", key, value)
}
func (e EmptyLogger) Version(v Version) {
	fmt.Println("Version", v)
}
func (e EmptyLogger) Info(msg string, args ...interface{}) {
	fmt.Println(msg, args)
}
func (e EmptyLogger) DownloadStarted(url string) {
	fmt.Println("DownloadStarted", url)
}
func (e EmptyLogger) DownloadFinished(url string) {
	fmt.Println("DownloadFinished", url)
}
func (e EmptyLogger) UploadStarted(url string) {
	fmt.Println("DownloadStarted", url)
}
func (e EmptyLogger) UploadFinished(url string) {
	fmt.Println("DownloadFinished", url)
}
func (e EmptyLogger) Running() {
	fmt.Println("Running")
}
func (e EmptyLogger) Stdout() io.Writer {
	return nil
}
func (e EmptyLogger) Stderr() io.Writer {
	return nil
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
