package main

import (
  "context"
  "fmt"
  tug "github.com/buchanae/tugboat"
  "github.com/buchanae/tugboat/docker"
  "github.com/buchanae/tugboat/storage/local"
)

func main() {
  ctx := context.Background()
  log := tug.EmptyLogger{}
  store := &local.Local{}

	stage, err := tug.NewStage("foo", 0755)
  if err != nil {
    panic(err)
  }
  //stage.LeaveDir = true

  exec := &docker.Docker{
    Logger: log,
  }

  task := &tug.Task{
    ID: "test1",
    ContainerImage: "alpine",
    Command: []string{"echo", "hello tugboat!"},
    Stdout: "out.txt",
    Outputs: []tug.File{
      {
        URL: "output",
        Path: "out.txt",
      },
    },
  }

  err = tug.Run(ctx, task, stage, log, store, exec)
  if err != nil {
    fmt.Println("RESULT", err)
  } else {
    fmt.Println("Success")
  }
}
